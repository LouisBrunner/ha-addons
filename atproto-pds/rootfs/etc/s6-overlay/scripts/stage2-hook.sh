#!/command/with-contenv bashio
# shellcheck shell=bash
set -e

source /etc/s6-overlay/scripts/logging.sh

if bashio::config.true 'gatekeeper.enabled'; then
  bashio::log.info "Gatekeeper enabled, adding caddy and gatekeeper services"
  touch /etc/s6-overlay/s6-rc.d/user/contents.d/caddy
  touch /etc/s6-overlay/s6-rc.d/user/contents.d/gatekeeper
fi
