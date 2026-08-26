#!/command/with-contenv bashio
# shellcheck shell=bash

bashio::config.require 'hostname'
HOSTNAME="$(bashio::config 'hostname')"
if [[ "${HOSTNAME}" == *"://"* || "${HOSTNAME}" == *"/"* ]]; then
	bashio::log.fatal "hostname must not include a scheme or path, e.g. 'app.mydomain.com' not 'https://app.mydomain.com/'"
	exit 1
fi
