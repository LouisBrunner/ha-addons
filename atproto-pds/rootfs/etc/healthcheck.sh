#!/command/with-contenv bashio
# shellcheck shell=bash
set -e

curl -fsS -o /dev/null http://127.0.0.1:3001/xrpc/_health
nc -z 127.0.0.1 3000

if bashio::config.true 'gatekeeper.enabled'; then
	nc -z 127.0.0.1 3002
fi
