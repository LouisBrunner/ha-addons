#!/command/with-contenv bashio
# shellcheck shell=bash
set -e

source /etc/s6-overlay/scripts/logging.sh

if bashio::config.true 'gatekeeper.enabled'; then
	bashio::log.info "Gatekeeper enabled, launch service"
	touch /etc/s6-overlay/s6-rc.d/user/contents.d/gatekeeper
fi
