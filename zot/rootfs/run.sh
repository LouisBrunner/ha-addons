#!/command/with-contenv bashio
# shellcheck shell=bash
set -e

bashio::config.require 'oidc.issuer'
bashio::config.require 'oidc.client_id'
bashio::config.require.safe_password 'oidc.client_secret'
bashio::config.require 'external_url'

mkdir -p /data/registry

SESSION_KEYS_FILE=/data/session-keys.json
if ! bashio::fs.file_exists "${SESSION_KEYS_FILE}"; then
	bashio::log.info "First run: generating session keys..."
	hash_key="$(openssl rand -hex 16)"
	encrypt_key="$(openssl rand -hex 16)"
	jq -n --arg hashKey "${hash_key}" --arg encryptKey "${encrypt_key}" \
		'{"hashKey": $hashKey, "encryptKey": $encryptKey}' >"${SESSION_KEYS_FILE}"
	chmod 600 "${SESSION_KEYS_FILE}"
fi

if bashio::config.has_value 'sync'; then
	bashio::log.info "Generating sync-credentials.json..."
	bashio::addon.config | tempio -template /etc/zot/sync-credentials.json.gtpl -out /etc/zot/sync-credentials.json
	chmod 600 /etc/zot/sync-credentials.json
fi

bashio::log.info "Generating config.json..."
bashio::addon.config | tempio -template /etc/zot/config.json.gtpl -out /etc/zot/config.json
chmod 600 /etc/zot/config.json
bashio::log.debug "config.json: $(jq '.http.auth.openid.providers.oidc.clientsecret = "***REDACTED***"' /etc/zot/config.json)"

bashio::log.info "Starting Zot on :5000..."
if bashio::debug; then
	bashio::log.debug "Environment:"
	while IFS='=' read -r name value; do
		case "${name^^}" in
		*SECRET* | *PASSWORD*) echo "${name}=***REDACTED***" ;;
		*) echo "${name}=${value}" ;;
		esac
	done < <(env)
fi
exec zot serve /etc/zot/config.json
