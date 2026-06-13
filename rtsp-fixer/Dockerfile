ARG BUILD_FROM=stop-complaining

FROM golang:1.26-alpine AS go-builder
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
RUN apk add --no-cache openh264
COPY --from=go-build /code/server /server
COPY run.sh /
RUN chmod a+x /run.sh
CMD ["/run.sh"]
