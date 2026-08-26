FROM golang:1.26.4-alpine3.24 AS go-builder
RUN apk add --no-cache gcc~15.2.0 musl-dev~1.2.6
WORKDIR /code
RUN go env

FROM go-builder AS go-deps
COPY server/go.* .
RUN --mount=type=cache,target=/go/pkg/mod \
  go mod download

FROM go-deps AS go-build
COPY server .
RUN --mount=type=cache,target=/root/.cache/go-build \
  --mount=type=cache,target=/go/pkg/mod \
  go build ./cmd/server

FROM ghcr.io/home-assistant/base:3.23-2026.06.1
RUN apk add --no-cache openh264~2.6.0 netcat-openbsd~1.234.1
COPY --from=go-build /code/server /server
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
  org.opencontainers.image.documentation="https://github.com/LouisBrunner/ha-addons/blob/main/rtsp-fixer/README.md" \
  org.opencontainers.image.version="${BUILD_VERSION}"

HEALTHCHECK --interval=30s --timeout=3s CMD nc -z 127.0.0.1 6667 || exit 1
