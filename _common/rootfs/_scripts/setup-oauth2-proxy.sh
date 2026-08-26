#!/bin/sh
set -eu

TARGETARCH="$1"
OAUTH2_PROXY_VERSION="$2"

case "${TARGETARCH}" in
	amd64) SHA256="0ae5a43adde4d6c5081ba018e70a76041f496377b12a173da36b419082dd1ab6" ;;
	arm64) SHA256="62452322a71e958d4d6911f799bc07921212a5f3bc45e39b63746e422d52ea33" ;;
	*)
		echo "unsupported arch: ${TARGETARCH}" >&2
		exit 1
		;;
esac

curl -fsSL -o /tmp/oauth2-proxy.tar.gz "https://github.com/oauth2-proxy/oauth2-proxy/releases/download/${OAUTH2_PROXY_VERSION}/oauth2-proxy-${OAUTH2_PROXY_VERSION}.linux-${TARGETARCH}.tar.gz"
echo "${SHA256}  /tmp/oauth2-proxy.tar.gz" | sha256sum -c -
tar -xzf /tmp/oauth2-proxy.tar.gz -C /usr/bin --strip-components=1 \
	"oauth2-proxy-${OAUTH2_PROXY_VERSION}.linux-${TARGETARCH}/oauth2-proxy"
chmod +x /usr/bin/oauth2-proxy
rm /tmp/oauth2-proxy.tar.gz
