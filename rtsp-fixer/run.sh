#!/usr/bin/with-contenv bashio
# shellcheck shell=bash
set -e

bashio::log.info "Starting RTSP Stream Fixer..."

export STREAMS=$(bashio::config 'streams')
export DEBUG=$(bashio::config 'debug')
export PORT=6666

exec /server
