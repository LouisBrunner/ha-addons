#!/command/with-contenv bashio
# shellcheck shell=bash

source /bin/validate-hostname.sh

CADDY_EXTRA_VARS="$(jq -n --arg port "8080" '{upstream_port: $port}')"
export CADDY_EXTRA_VARS
