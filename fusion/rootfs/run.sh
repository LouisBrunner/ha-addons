#!/command/with-contenv bashio
# shellcheck shell=bash
set -e

bashio::config.require 'external_url'
bashio::config.require 'oidc.issuer'
bashio::config.require 'oidc.client_id'
bashio::config.require.safe_password 'oidc.client_secret'

EXTERNAL_URL="$(bashio::config 'external_url')"
if [[ "${EXTERNAL_URL}" != http://* && "${EXTERNAL_URL}" != https://* ]]; then
	bashio::log.fatal "external_url must include a scheme, e.g. 'https://fusion.mydomain.com'"
	exit 1
fi

mkdir -p /data

PASSWORD_FILE=/data/fusion_password
if ! bashio::fs.file_exists "${PASSWORD_FILE}"; then
	bashio::log.info "First run: generating a random password (OIDC is the intended login path)..."
	openssl rand -base64 32 | tr -d '\n' >"${PASSWORD_FILE}"
	chmod 600 "${PASSWORD_FILE}"
fi

export FUSION_PASSWORD
FUSION_PASSWORD="$(cat "${PASSWORD_FILE}")"
export FUSION_DB_PATH=/data/fusion.db
export FUSION_PORT=8080

export FUSION_OIDC_ISSUER
FUSION_OIDC_ISSUER="$(bashio::config 'oidc.issuer')"
export FUSION_OIDC_CLIENT_ID
FUSION_OIDC_CLIENT_ID="$(bashio::config 'oidc.client_id')"
export FUSION_OIDC_CLIENT_SECRET
FUSION_OIDC_CLIENT_SECRET="$(bashio::config 'oidc.client_secret')"
export FUSION_OIDC_REDIRECT_URI="${EXTERNAL_URL}/api/oidc/callback"

if bashio::config.has_value 'oidc.allowed_user'; then
	export FUSION_OIDC_ALLOWED_USER
	FUSION_OIDC_ALLOWED_USER="$(bashio::config 'oidc.allowed_user')"
fi

if bashio::config.has_value 'fever_username'; then
	export FUSION_FEVER_USERNAME
	FUSION_FEVER_USERNAME="$(bashio::config 'fever_username')"
fi

if bashio::config.true 'allow_private_feeds'; then
	export FUSION_ALLOW_PRIVATE_FEEDS=true
fi

if bashio::config.true 'debug'; then
	export FUSION_LOG_LEVEL=DEBUG
fi

bashio::log.info "Starting Fusion on 0.0.0.0:8080..."
cd /data
exec fusion
