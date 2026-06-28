#!/command/with-contenv bashio
# shellcheck shell=bash
set -e

bashio::config.require 'hostname'
PDS_HOSTNAME="$(bashio::config 'hostname')"
export PDS_HOSTNAME

bashio::config.require 'smtp_url'
PDS_EMAIL_SMTP_URL="$(bashio::config 'smtp_url')"
export PDS_EMAIL_SMTP_URL
PDS_EMAIL_FROM_ADDRESS="$(bashio::config 'smtp_from' "no-reply@${PDS_HOSTNAME}")"
export PDS_EMAIL_FROM_ADDRESS

export PDS_DATA_DIRECTORY=/data/pds
export PDS_PUBLIC_URL="https://${PDS_HOSTNAME}"
