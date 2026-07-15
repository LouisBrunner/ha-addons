# syntax=docker/dockerfile:1.17

# FROM fatfingers23/pds_gatekeeper:0.1.0.7 AS gatekeeper # FIXME: can't use, glibc
FROM rust:1.96.1-alpine3.22 AS gatekeeper
RUN apk add --no-cache musl-dev pkgconf openssl-dev openssl-libs-static
WORKDIR /app
#checkov:skip=CKV_DOCKER_4: same behavior as COPY + extraction
ADD --unpack https://tangled.org/did:plc:p3aqxmyn6f3pkm4dcl5co5gp/archive/eef639e04811b734f426548a63fc8911886c3b48?format=tar.gz /src
RUN cp -a /src/*/. /app/
RUN --mount=type=cache,target=/usr/local/cargo/registry \
  --mount=type=cache,target=/usr/local/cargo/git/db \
  --mount=type=cache,target=/app/target \
  cargo build --release && cp target/release/pds_gatekeeper /pds_gatekeeper

FROM ghcr.io/bluesky-social/pds:0.4.5009 AS pds

FROM ghcr.io/home-assistant/base:3.24
RUN apk add --no-cache caddy nodejs-current openssl curl~8.19.0 netcat-openbsd~1.234.1

ENV S6_STAGE2_HOOK=/etc/s6-overlay/scripts/stage2-hook.sh

COPY --from=gatekeeper /pds_gatekeeper /usr/local/bin/pds_gatekeeper
COPY --from=pds /app /usr/local/bin/pds
COPY --from=pds /usr/local/bin/goat /usr/local/bin/goat
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
  org.opencontainers.image.documentation="https://github.com/${BUILD_REPOSITORY}/blob/main/atproton-pds/README.md" \
  org.opencontainers.image.created=${BUILD_DATE} \
  org.opencontainers.image.revision=${BUILD_REF} \
  org.opencontainers.image.version=${BUILD_VERSION}

HEALTHCHECK --interval=30s --timeout=5s --start-period=30s CMD /etc/healthcheck.sh
