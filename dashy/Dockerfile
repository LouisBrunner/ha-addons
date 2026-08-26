# syntax=docker/dockerfile:1
FROM ghcr.io/lissy93/dashy:4.5.12 AS app
FROM oven/bun:1-alpine AS bun

FROM ghcr.io/home-assistant/base:3.23
RUN apk add --no-cache curl~8.20.0 yq-go~4.49.2
COPY --from=bun /usr/local/bin/bun /usr/local/bin/bun
COPY --from=app /app /app
COPY rootfs /
CMD ["/run.sh"]

ARG BUILD_VERSION
LABEL \
  maintainer="Louis Brunner <opensource@e.lbrunner.net>" \
  org.opencontainers.image.authors="Louis Brunner <opensource@e.lbrunner.net>" \
  org.opencontainers.image.vendor="Louis Brunner's Home Assistant Apps" \
  org.opencontainers.image.licenses="MPL-2.0" \
  org.opencontainers.image.source="https://github.com/LouisBrunner/ha-addons" \
  org.opencontainers.image.url="https://github.com/LouisBrunner/ha-addons" \
  org.opencontainers.image.documentation="https://github.com/LouisBrunner/ha-addons/blob/main/dashy/README.md" \
  org.opencontainers.image.version="${BUILD_VERSION}"

HEALTHCHECK --interval=5m --timeout=10s --start-period=20s --retries=3 \
  CMD bun /app/services/healthcheck.js
