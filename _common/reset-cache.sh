#!/bin/bash
set -euo pipefail

slug="local_${1:?usage: reset-cache.sh <target>}"

jq --arg slug "${slug}" '.user[$slug].version = "0.0.0"' /mnt/supervisor/apps.json >/tmp/apps.json.new
mv /tmp/apps.json.new /mnt/supervisor/apps.json

ha supervisor restart >/dev/null

for _ in $(seq 1 10); do
	ha apps update "${slug}" >/dev/null 2>&1 && exit 0
	sleep 1
done

echo "reset-cache: gave up waiting for 'ha apps update ${slug}' to succeed" >&2
exit 1
