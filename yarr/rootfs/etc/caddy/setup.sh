#!/command/with-contenv bashio
# shellcheck shell=bash

source /bin/validate-hostname.sh

if bashio::config.has_value 'fever_password'; then
	fever_username="fever"
	fever_password="$(bashio::config 'fever_password')"
	# Mirrors yarr's own auth cookie derivation (src/server/auth/auth.go, function secret):
	# hex(HMAC-SHA256(key=password, msg=username)) — yarr's web UI only ever checks this cookie,
	# it has no Basic Auth support, so this is the only way to bridge OIDC into its session model.
	auth_cookie_hmac="$(printf '%s' "${fever_username}" | openssl dgst -sha256 -hmac "${fever_password}" -r | awk '{print $1}')"

	mkdir -p /etc/caddy/extra-routes.d
	cp /etc/caddy/fever-route.caddy.disabled /etc/caddy/extra-routes.d/fever.caddy

	jq -n --arg cookie "${fever_username}:${auth_cookie_hmac}" '{cookie: $cookie}' \
		| tempio -template /etc/caddy/extra-headers.d/fever-cookie.caddy.gtpl -out /etc/caddy/extra-headers.d/fever-cookie.caddy
fi

CADDY_EXTRA_VARS="$(jq -n --arg port "7070" '{upstream_port: $port}')"
export CADDY_EXTRA_VARS
