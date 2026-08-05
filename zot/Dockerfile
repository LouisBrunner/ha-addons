# syntax=docker/dockerfile:1
FROM ghcr.io/project-zot/zot:v2.1.18 AS zot

FROM ghcr.io/home-assistant/base:3.23
RUN apk add --no-cache jq~1.8.1 curl~8.19.0 gcompat~1.1.0 openssl~3.5.7

ARG TARGETARCH
COPY --from=zot /usr/local/bin/zot-linux-${TARGETARCH} /usr/bin/zot
COPY rootfs /
CMD ["/run.sh"]

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
  org.opencontainers.image.documentation="https://github.com/${BUILD_REPOSITORY}/blob/main/zot/README.md" \
  org.opencontainers.image.created=${BUILD_DATE} \
  org.opencontainers.image.revision=${BUILD_REF} \
  org.opencontainers.image.version=${BUILD_VERSION}

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s \
  CMD status="$(curl -sS -o /dev/null -w '%{http_code}' http://127.0.0.1:5000/v2/)"; [ "${status}" = "200" ] || [ "${status}" = "401" ]
