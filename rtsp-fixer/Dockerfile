ARG BUILD_FROM=ghcr.io/home-assistant/base:3.23

FROM golang:1.26.4-alpine3.24 AS go-builder
RUN apk add --no-cache gcc==15.2.0-r5 musl-dev==1.2.6-r24
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

FROM $BUILD_FROM
RUN apk add --no-cache openh264==2.6.0-r0
COPY --from=go-build /code/server /server
CMD ["/run.sh"]
