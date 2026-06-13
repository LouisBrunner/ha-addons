STREAMS ?= '{"name":"test","url":"rtsp://192.168.0.66:554/onvif1","fix_force_tcp_in_transport":true}'
PORT ?= 6668
PORT_THUMBNAILS ?= 6669
DEBUG ?= true
BASE_FOLDER = "./data"

run:
	STREAMS=$(STREAMS) PORT=$(PORT) PORT_THUMBNAILS=$(PORT_THUMBNAILS) DEBUG=$(DEBUG) BASE_FOLDER=$(BASE_FOLDER) go run ./cmd/server
.PHONY: run

vlc-test:
	vlc --rtsp-tcp --verbose 4 rtsp://:@localhost:6668/test
.PHONY: vlc-test

mod-tidy:
	go mod tidy
.PHONY: mod-tidy
