# syntax=docker/dockerfile:1

FROM scratch AS common
#checkov:skip=CKV_DOCKER_4: same behavior as COPY + extraction
ADD --unpack http://localhost:8787/common.tar.gz /

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

RUN --mount=from=common,source=_scripts/setup-oauth2-proxy.sh,target=/setup-oauth2-proxy.sh \
  sh /setup-oauth2-proxy.sh "${TARGETARCH}" "${OAUTH2_PROXY_VERSION}"

FROM ghcr.io/home-assistant/base:3.23
RUN apk add --no-cache caddy~2.11.4 openssl~3.5.7 jq~1.8.1
COPY --from=builder /usr/bin/yarr /usr/bin/yarr
COPY --from=builder /usr/bin/oauth2-proxy /usr/bin/oauth2-proxy
COPY --from=common /base /
COPY --from=common /caddy /
COPY --from=common /oauth2-proxy /
COPY rootfs /

ARG BUILD_VERSION
LABEL \
  maintainer="Louis Brunner <opensource@e.lbrunner.net>" \
  org.opencontainers.image.authors="Louis Brunner <opensource@e.lbrunner.net>" \
  org.opencontainers.image.vendor="Louis Brunner's Home Assistant Apps" \
  org.opencontainers.image.licenses="MPL-2.0" \
  org.opencontainers.image.source="https://github.com/LouisBrunner/ha-addons" \
  org.opencontainers.image.url="https://github.com/LouisBrunner/ha-addons" \
  org.opencontainers.image.documentation="https://github.com/LouisBrunner/ha-addons/blob/main/yarr/README.md" \
  org.opencontainers.image.version="${BUILD_VERSION}"

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s \
  CMD /bin/caddy-healthcheck.sh
