# syntax=docker/dockerfile:1
FROM ghcr.io/project-zot/zot:v2.1.18 AS zot

FROM ghcr.io/home-assistant/base:3.23
RUN apk add --no-cache jq~1.8.1 curl~8.19.0 gcompat~1.1.0 openssl~3.5.7 netcat-openbsd~1.234.1

ARG TARGETARCH
COPY --from=zot /usr/local/bin/zot-linux-${TARGETARCH} /usr/bin/zot
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
  org.opencontainers.image.documentation="https://github.com/LouisBrunner/ha-addons/blob/main/zot/README.md" \
  org.opencontainers.image.version="${BUILD_VERSION}"

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s \
  CMD nc -z 127.0.0.1 5000
