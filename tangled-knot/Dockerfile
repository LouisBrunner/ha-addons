FROM atcr.io/tangled.org/knot:v1.15.0-alpha
ARG BASHIO_VERSION=0.18.1

RUN apk add --no-cache bash curl jq \
    && mkdir -p /usr/lib/bashio \
    && curl -fsSL "https://github.com/hassio-addons/bashio/archive/v${BASHIO_VERSION}.tar.gz" \
    | tar -xzm --strip-components=2 -C /usr/lib/bashio

COPY run.sh /
RUN chmod a+x /run.sh
CMD ["/run.sh"]
