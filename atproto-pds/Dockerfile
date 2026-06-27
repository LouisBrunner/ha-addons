FROM ghcr.io/bluesky-social/pds:0.4.5009
ARG BASHIO_VERSION=v0.18.1

RUN apk add --no-cache curl jq openssl bash \
    && mkdir -p /usr/lib/bashio \
    && curl -fsSL "https://github.com/hassio-addons/bashio/archive/${BASHIO_VERSION}.tar.gz" \
    | tar -xzm --strip-components=2 -C /usr/lib/bashio && ln -s /usr/lib/bashio/bashio /usr/bin/bashio

COPY --chmod=0755 run.sh /
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
  org.opencontainers.image.vendor="Louis Brunner " \
  org.opencontainers.image.authors="Louis Brunner <opensource@e.lbrunner.net>" \
  org.opencontainers.image.licenses="MPL-2.0" \
  org.opencontainers.image.url="https://github.com/${BUILD_REPOSITORY}" \
  org.opencontainers.image.source="https://github.com/${BUILD_REPOSITORY}" \
  org.opencontainers.image.documentation="https://github.com/${BUILD_REPOSITORY}/blob/main/atproton-pds/README.md" \
  org.opencontainers.image.created=${BUILD_DATE} \
  org.opencontainers.image.revision=${BUILD_REF} \
  org.opencontainers.image.version=${BUILD_VERSION}
