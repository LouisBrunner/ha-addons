# syntax=docker/dockerfile:1.17

FROM scratch AS common
#checkov:skip=CKV_DOCKER_4: same behavior as COPY + extraction
ADD --unpack http://localhost:8787/common.tar.gz /

# FROM fatfingers23/pds_gatekeeper:0.1.0.7 AS gatekeeper # FIXME: can't use, glibc
FROM rust:1.89.0-alpine3.22 AS gatekeeper
RUN apk add --no-cache musl-dev~1.2.5 pkgconf~2.4.3 openssl-dev~3.5.7 openssl-libs-static~3.5.7
WORKDIR /app
#checkov:skip=CKV_DOCKER_4: same behavior as COPY + extraction
ADD --unpack https://tangled.org/did:plc:p3aqxmyn6f3pkm4dcl5co5gp/archive/eef639e04811b734f426548a63fc8911886c3b48?format=tar.gz /src
RUN cp -a /src/*/. /app/
RUN --mount=type=cache,target=/usr/local/cargo/registry \
  --mount=type=cache,target=/usr/local/cargo/git/db \
  --mount=type=cache,target=/app/target \
  cargo build --release && cp target/release/pds_gatekeeper /pds_gatekeeper

FROM ghcr.io/bluesky-social/pds:0.4.5009 AS pds

FROM ghcr.io/home-assistant/base:3.23
RUN apk add --no-cache caddy~2.11.4 nodejs-current~24.15.0 openssl~3.5.7 curl~8.20.0 netcat-openbsd~1.234.1

ENV S6_STAGE2_HOOK=/etc/s6-overlay/scripts/stage2-hook.sh

COPY --from=gatekeeper /pds_gatekeeper /usr/local/bin/pds_gatekeeper
COPY --from=pds /app /usr/local/bin/pds
COPY --from=pds /usr/local/bin/goat /usr/local/bin/goat
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
  org.opencontainers.image.documentation="https://github.com/LouisBrunner/ha-addons/blob/main/atproto-pds/README.md" \
  org.opencontainers.image.version="${BUILD_VERSION}"

HEALTHCHECK --interval=30s --timeout=5s --start-period=30s CMD /etc/healthcheck.sh
