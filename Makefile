INSIDE_DOCKER = $(shell stat /.indocker 2>&1 >/dev/null && echo 1 || echo 0)
ifeq ($(INSIDE_DOCKER),1)
endif

all:
.PHONY: all

setup:
ifeq ($(INSIDE_DOCKER),1)
else
endif
.PHONY: setup

TARGET ?= unknown

deploy-local:
	scp -rP 2020 $(TARGET) root@homeassistant.local:/homeassistant/apps
.PHONY: deploy-local

devcontainer-start:
	@bash -c 'docker compose up --build --force-recreate -d && trap "docker compose down" EXIT INT TERM HUP && make logs'
.PHONY: devcontainer-start

devcontainer:
	docker compose exec -it devcontainer bash
.PHONY: devcontainer

logs:
ifeq ($(INSIDE_DOCKER),1)
	journalctl --no-tail -f -u hassio-supervisor -u hassio-bootstrap -u hassio-apparmor -u docker
else
	docker compose exec devcontainer make logs
endif
.PHONY: dev
