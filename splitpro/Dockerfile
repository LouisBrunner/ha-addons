# syntax=docker/dockerfile:1
FROM ossapps/splitpro:v2.1.5 AS app

FROM ghcr.io/home-assistant/base:3.23
RUN apk add --no-cache nodejs-current~24.15.0 postgresql18~18.4 postgresql-pg_cron~1.6.7 netcat-openbsd~1.234.1 openssl~3.5.7
COPY --from=app /app /app
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
  org.opencontainers.image.documentation="https://github.com/${BUILD_REPOSITORY}/blob/main/splitpro/README.md" \
  org.opencontainers.image.created=${BUILD_DATE} \
  org.opencontainers.image.revision=${BUILD_REF} \
  org.opencontainers.image.version=${BUILD_VERSION}

HEALTHCHECK --interval=30s --timeout=5s --start-period=30s \
  CMD nc -z 127.0.0.1 5432 && nc -z 127.0.0.1 3000
