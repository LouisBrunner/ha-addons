#!/usr/bin/with-contenv bashio
# shellcheck shell=bash
set -e

bashio::log.info "Starting RTSP Stream Fixer..."

STREAMS=$(bashio::config 'streams')
export STREAMS
DEBUG=$(bashio::config 'debug')
export DEBUG
export PORT=6666
export PORT_THUMBNAILS=6667

exec /server
