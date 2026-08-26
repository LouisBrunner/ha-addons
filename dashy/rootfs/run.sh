#!/command/with-contenv bashio
# shellcheck shell=bash
set -e

bashio::config.require 'external_url'
bashio::config.require 'oidc.issuer'
bashio::config.require 'oidc.client_id'

EXTERNAL_URL="$(bashio::config 'external_url')"
if [[ "${EXTERNAL_URL}" != http://* && "${EXTERNAL_URL}" != https://* ]]; then
	bashio::log.fatal "external_url must include a scheme, e.g. 'https://dashy.mydomain.com'"
	exit 1
fi

mkdir -p /data
rm -rf /app/user-data
ln -sfn /data /app/user-data

if ! bashio::fs.file_exists /data/conf.yml; then
	bashio::log.info "First run: seeding conf.yml..."
	cp /etc/dashy/conf.yml.default /data/conf.yml
fi

bashio::log.info "Applying OIDC options to conf.yml..."
bashio::addon.config | tempio -template /etc/dashy/oidc.yml.gtpl -out /tmp/oidc.yml
yq eval-all 'select(fileIndex==0) * select(fileIndex==1)' /data/conf.yml /tmp/oidc.yml > /tmp/conf.merged.yml
mv /tmp/conf.merged.yml /data/conf.yml

bashio::log.info "Starting Dashy on 0.0.0.0:8080..."
cd /app
exec bun server.js
