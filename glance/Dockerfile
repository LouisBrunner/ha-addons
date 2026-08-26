# syntax=docker/dockerfile:1

FROM scratch AS common
#checkov:skip=CKV_DOCKER_4: same behavior as COPY + extraction
ADD --unpack http://localhost:8787/common.tar.gz /

FROM glanceapp/glance:v0.8.5 AS app

FROM ghcr.io/home-assistant/base:3.23
RUN apk add --no-cache caddy~2.11.4 openssl~3.5.7 jq~1.8.1 curl~8.20.0
ARG OAUTH2_PROXY_VERSION=v7.15.3
ARG TARGETARCH
RUN --mount=from=common,source=_scripts/setup-oauth2-proxy.sh,target=/setup-oauth2-proxy.sh \
  sh /setup-oauth2-proxy.sh "${TARGETARCH}" "${OAUTH2_PROXY_VERSION}"

COPY --from=app /app/glance /usr/bin/glance
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
  org.opencontainers.image.documentation="https://github.com/LouisBrunner/ha-addons/blob/main/glance/README.md" \
  org.opencontainers.image.version="${BUILD_VERSION}"

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s \
  CMD /bin/caddy-healthcheck.sh
