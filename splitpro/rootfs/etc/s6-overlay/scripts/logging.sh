#!/command/with-contenv bashio
# shellcheck shell=bash
set -e

if bashio::config.true 'debug'; then
	export __BASHIO_LOG_LEVEL=${__BASHIO_LOG_LEVEL_DEBUG}
fi

print_env() {
	if bashio::debug; then
		bashio::log.debug "Environment:"
		while IFS='=' read -r name value; do
			case "${name^^}" in
			*SECRET* | *PASSWORD*) echo "${name}=***REDACTED***" ;;
			*) echo "${name}=${value}" ;;
			esac
		done < <(env)
	fi
}
