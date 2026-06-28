#!/usr/bin/env bashio
# shellcheck shell=bash
set -e

# Persist SSH host keys across restarts
mkdir -p /data/keys
rm -rf /etc/ssh/keys
ln -sf /data/keys /etc/ssh/keys

KNOT_SERVER_HOSTNAME="$(bashio::config 'hostname')"
export KNOT_SERVER_HOSTNAME
KNOT_SERVER_OWNER="$(bashio::config 'owner_did')"
export KNOT_SERVER_OWNER
export KNOT_SERVER_DB_PATH=/data/knotserver.db
export KNOT_REPO_SCAN_PATH=/data/repositories
export KNOT_SERVER_INTERNAL_LISTEN_ADDR=0.0.0.0:5555

mkdir -p "${KNOT_REPO_SCAN_PATH}"

bashio::log.info "Starting Tangled knot on ${KNOT_SERVER_HOSTNAME}..."
if bashio::debug; then
	bashio::log.debug "Environment:"
	env
fi
exec /init
