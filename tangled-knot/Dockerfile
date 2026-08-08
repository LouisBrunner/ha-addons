FROM atcr.io/tangled.org/knot:v1.15.0-alpha
ARG BASHIO_VERSION=v0.18.1

SHELL ["/bin/ash", "-o", "pipefail", "-c"]
RUN apk add --no-cache curl~8.21.0 jq~1.8.1 openssl~3.5.7 bash~5.3.9 netcat-openbsd~1.234.1 \
  && mkdir -p /usr/lib/bashio \
  && curl -fsSL "https://github.com/hassio-addons/bashio/archive/${BASHIO_VERSION}.tar.gz" \
  | tar -xzm --strip-components=2 -C /usr/lib/bashio && ln -s /usr/lib/bashio/bashio /usr/bin/bashio

CMD ["/run.sh"]

HEALTHCHECK --interval=30s --timeout=5s --start-period=30s \
  CMD nc -z 127.0.0.1 5555

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
  org.opencontainers.image.documentation="https://github.com/${BUILD_REPOSITORY}/blob/main/tangled-knot/README.md" \
  org.opencontainers.image.created=${BUILD_DATE} \
  org.opencontainers.image.revision=${BUILD_REF} \
  org.opencontainers.image.version=${BUILD_VERSION}
