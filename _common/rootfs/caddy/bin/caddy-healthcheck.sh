#!/bin/sh
set -eu
exec wget -q -O /dev/null "http://127.0.0.1:${CADDY_HEALTHCHECK_PORT:-3000}/healthz"
