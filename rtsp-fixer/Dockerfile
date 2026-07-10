FROM golang:1.26.4-alpine3.24 AS go-builder
RUN apk add --no-cache gcc==15.2.0-r5 musl-dev==1.2.6-r2
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
RUN apk add --no-cache openh264==2.6.0-r0 netcat-openbsd==1.234.1-r0
COPY --from=go-build /code/server /server
CMD ["/run.sh"]

HEALTHCHECK --interval=30s --timeout=3s CMD nc -z 127.0.0.1 6667 || exit 1
