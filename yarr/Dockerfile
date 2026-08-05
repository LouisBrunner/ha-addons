# syntax=docker/dockerfile:1
FROM ghcr.io/home-assistant/base:3.23 AS builder
RUN apk add --no-cache curl~8.19.0 unzip~6.0 tar~1.35
ARG YARR_VERSION=v2.8
ARG OAUTH2_PROXY_VERSION=v7.15.3
ARG TARGETARCH

RUN set -eux; \
  case "${TARGETARCH}" in \
  amd64) SHA256="2b8f7ffdcbf939a77192cf69551a9cfbb4ba9fb3f2be7c6ff75a6afb94168bbd" ;; \
  arm64) SHA256="6fbd40f437522367bcb07139d00d77351efce8fd8aee6c3cc9f25efe290234b0" ;; \
  *) echo "unsupported arch: ${TARGETARCH}" >&2; exit 1 ;; \
  esac; \
  curl -fsSL -o /tmp/yarr.zip "https://github.com/nkanaev/yarr/releases/download/${YARR_VERSION}/yarr_linux_${TARGETARCH}.zip"; \
  echo "${SHA256}  /tmp/yarr.zip" | sha256sum -c -; \
  unzip -q /tmp/yarr.zip -d /usr/bin; \
  chmod +x /usr/bin/yarr; \
  rm /tmp/yarr.zip

RUN set -eux; \
  case "${TARGETARCH}" in \
  amd64) SHA256="0ae5a43adde4d6c5081ba018e70a76041f496377b12a173da36b419082dd1ab6" ;; \
  arm64) SHA256="62452322a71e958d4d6911f799bc07921212a5f3bc45e39b63746e422d52ea33" ;; \
  *) echo "unsupported arch: ${TARGETARCH}" >&2; exit 1 ;; \
  esac; \
  curl -fsSL -o /tmp/oauth2-proxy.tar.gz "https://github.com/oauth2-proxy/oauth2-proxy/releases/download/${OAUTH2_PROXY_VERSION}/oauth2-proxy-${OAUTH2_PROXY_VERSION}.linux-${TARGETARCH}.tar.gz"; \
  echo "${SHA256}  /tmp/oauth2-proxy.tar.gz" | sha256sum -c -; \
  tar -xzf /tmp/oauth2-proxy.tar.gz -C /usr/bin --strip-components=1 \
    "oauth2-proxy-${OAUTH2_PROXY_VERSION}.linux-${TARGETARCH}/oauth2-proxy"; \
  chmod +x /usr/bin/oauth2-proxy; \
  rm /tmp/oauth2-proxy.tar.gz

FROM ghcr.io/home-assistant/base:3.23
RUN apk add --no-cache caddy~2.11.4 openssl~3.5.7 jq~1.8.1
COPY --from=builder /usr/bin/yarr /usr/bin/yarr
COPY --from=builder /usr/bin/oauth2-proxy /usr/bin/oauth2-proxy
COPY rootfs /

ARG BUILD_ARCH
ARG BUILD_DATE
ARG BUILD_DESCRIPTION
ARG BUILD_NAME
ARG BUILD_REF
ARG BUILD_REPOSITORY
ARG BUILD_VERSION
LABEL \
  io.hass.type="app" \
  io.hass.name="${BUILD_NAME}" \
  io.hass.description="${BUILD_DESCRIPTION}" \
  io.hass.version="${BUILD_VERSION}" \
  io.hass.arch="${BUILD_ARCH}" \
  maintainer="Louis Brunner <opensource@e.lbrunner.net>" \
  org.opencontainers.image.title="${BUILD_NAME}" \
  org.opencontainers.image.description="${BUILD_DESCRIPTION}" \
  org.opencontainers.image.vendor="Louis Brunner's Home Assistant Apps" \
  org.opencontainers.image.authors="Louis Brunner <opensource@e.lbrunner.net>" \
  org.opencontainers.image.licenses="MPL-2.0" \
  org.opencontainers.image.url="https://github.com/${BUILD_REPOSITORY}" \
  org.opencontainers.image.source="https://github.com/${BUILD_REPOSITORY}" \
  org.opencontainers.image.documentation="https://github.com/${BUILD_REPOSITORY}/blob/main/yarr/README.md" \
  org.opencontainers.image.created=${BUILD_DATE} \
  org.opencontainers.image.revision=${BUILD_REF} \
  org.opencontainers.image.version=${BUILD_VERSION}

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s \
  CMD wget -q -O /dev/null http://127.0.0.1:3000/healthz || exit 1
