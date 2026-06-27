#!/usr/bin/env bashio
# shellcheck shell=bash
set -e

SECRETS_FILE=/data/secrets.env

if ! bashio::fs.file_exists "${SECRETS_FILE}"; then
  bashio::log.info "First run: generating secrets..."

  jwt_secret="$(openssl rand -hex 32)"
  plc_key="$(openssl ecparam --name secp256k1 --genkey --noout --outform DER | tail -c +8 | head -c 32 | xxd -p -c 32)"

  if [[ -z "${jwt_secret}" ]] || [[ -z "${plc_key}" ]]; then
    bashio::log.fatal "Failed to generate secrets — is openssl working?"
    exit 1
  fi

  printf "export PDS_JWT_SECRET=%s\nexport PDS_PLC_ROTATION_KEY_K256_PRIVATE_KEY_HEX=%s\n" "${jwt_secret}" "${plc_key}" >"${SECRETS_FILE}"
  bashio::log.info "Secrets written to ${SECRETS_FILE}"
fi

source "${SECRETS_FILE}"

bashio::config.require 'hostname'
bashio::config.require.safe_password "admin_password"

export PDS_HOSTNAME="$(bashio::config 'hostname')"
export PDS_ADMIN_PASSWORD="$(bashio::config 'admin_password')"
export PDS_INVITE_REQUIRED="$(bashio::config 'invite_required' 'true')"

export PDS_DATA_DIRECTORY=/data
export PDS_BLOBSTORE_DISK_LOCATION=/data/blobs
export PDS_PORT=3000
export LOG_ENABLED=true

# Federation — values match the official installer defaults
export PDS_BLOB_UPLOAD_LIMIT=104857600
export PDS_DID_PLC_URL=https://plc.directory
export PDS_BSKY_APP_VIEW_URL=https://api.bsky.app
export PDS_BSKY_APP_VIEW_DID=did:web:api.bsky.app
export PDS_REPORT_SERVICE_URL=https://mod.bsky.app
export PDS_REPORT_SERVICE_DID=did:plc:ar7c4by46qjdydhdevvrndac
export PDS_CRAWLERS=https://bsky.network
export PDS_RATE_LIMITS_ENABLED=true

bashio::log.info "Starting AT Protocol PDS on ${PDS_HOSTNAME}:${PDS_PORT}..."
if bashio::debug; then
  bashio::log.debug "Environment:"
  env
fi
exec node --enable-source-maps /app/index.ts
