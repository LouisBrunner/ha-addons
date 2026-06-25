FROM ghcr.io/bluesky-social/pds:v0.4.5006
ARG BASHIO_VERSION=0.18.1

RUN apt-get update \
    && apt-get install -y --no-install-recommends curl jq openssl \
    && rm -rf /var/lib/apt/lists/* \
    && mkdir -p /usr/lib/bashio \
    && curl -fsSL "https://github.com/hassio-addons/bashio/archive/v${BASHIO_VERSION}.tar.gz" \
    | tar -xzm --strip-components=2 -C /usr/lib/bashio

COPY run.sh /
RUN chmod a+x /run.sh
CMD ["/run.sh"]
