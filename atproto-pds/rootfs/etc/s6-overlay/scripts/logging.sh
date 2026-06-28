#!/command/with-contenv bashio
# shellcheck shell=bash
set -e

if bashio::config.true 'debug'; then
  export __BASHIO_LOG_LEVEL=${__BASHIO_LOG_LEVEL_DEBUG}
fi

print_env() {
  if bashio::debug; then
    bashio::log.debug "Environment:"
    env
  fi
}
