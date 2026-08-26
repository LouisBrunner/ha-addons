# syntax=docker/dockerfile:1

FROM scratch AS common
#checkov:skip=CKV_DOCKER_4: same behavior as COPY + extraction
ADD --unpack http://localhost:8787/common.tar.gz /

FROM caddy:2.11.4-builder-alpine AS builder
RUN xcaddy build \
  --with github.com/mholt/caddy-ratelimit@v0.1.0

FROM ghcr.io/home-assistant/base:3.23
RUN apk add --no-cache jq~1.8.1
COPY --from=builder /usr/bin/caddy /usr/bin/caddy
COPY --from=common /base /
COPY --from=common /caddy /
COPY rootfs /

ARG BUILD_VERSION
LABEL \
  maintainer="Louis Brunner <opensource@e.lbrunner.net>" \
  org.opencontainers.image.authors="Louis Brunner <opensource@e.lbrunner.net>" \
  org.opencontainers.image.vendor="Louis Brunner's Home Assistant Apps" \
  org.opencontainers.image.licenses="MPL-2.0" \
  org.opencontainers.image.source="https://github.com/LouisBrunner/ha-addons" \
  org.opencontainers.image.url="https://github.com/LouisBrunner/ha-addons" \
  org.opencontainers.image.documentation="https://github.com/LouisBrunner/ha-addons/blob/main/caddy/README.md" \
  org.opencontainers.image.version="${BUILD_VERSION}"

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s \
  CMD CADDY_HEALTHCHECK_PORT=8080 /bin/caddy-healthcheck.sh
