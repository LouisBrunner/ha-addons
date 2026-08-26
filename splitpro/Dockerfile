# syntax=docker/dockerfile:1

FROM scratch AS common
#checkov:skip=CKV_DOCKER_4: same behavior as COPY + extraction
ADD --unpack http://localhost:8787/common.tar.gz /

FROM ossapps/splitpro:v2.1.5 AS app

FROM ghcr.io/home-assistant/base:3.23
RUN apk add --no-cache nodejs-current~24.15.0 npm~11.11.0 postgresql18~18.6 postgresql-pg_cron~1.6.7 netcat-openbsd~1.234.1 openssl~3.5.7
COPY --from=app /app /app
RUN npm install -g prisma@6.19.3
COPY --from=common /base /
COPY rootfs /

ARG BUILD_VERSION
LABEL \
  maintainer="Louis Brunner <opensource@e.lbrunner.net>" \
  org.opencontainers.image.authors="Louis Brunner <opensource@e.lbrunner.net>" \
  org.opencontainers.image.vendor="Louis Brunner's Home Assistant Apps" \
  org.opencontainers.image.licenses="MPL-2.0" \
  org.opencontainers.image.source="https://github.com/LouisBrunner/ha-addons" \
  org.opencontainers.image.url="https://github.com/LouisBrunner/ha-addons" \
  org.opencontainers.image.documentation="https://github.com/LouisBrunner/ha-addons/blob/main/splitpro/README.md" \
  org.opencontainers.image.version="${BUILD_VERSION}"

HEALTHCHECK --interval=30s --timeout=5s --start-period=30s \
  CMD nc -z 127.0.0.1 5432 && nc -z 127.0.0.1 3000
