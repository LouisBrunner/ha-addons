TARGET ?= unknown

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

deploy-local:
	scp -rP 2020 $(TARGET) root@homeassistant.local:/homeassistant/apps
.PHONY: deploy-local

devcontainer-start:
	@bash -c 'docker compose up --build --force-recreate -d && trap "docker compose down" EXIT INT TERM HUP && make logs'
.PHONY: devcontainer-start

devcontainer:
	docker compose exec -it devcontainer bash
.PHONY: devcontainer

dev:
ifeq ($(INSIDE_DOCKER),1)
	@echo "> Unsupported inside the container"
else
	@echo "# Developing $(TARGET)"
	@make TARGET=$(TARGET) rebuild
	@echo "# Watching for changes..."
	@fswatch $(TARGET) | xargs -I{} make TARGET=$(TARGET) rebuild
endif
.PHONY: dev

rebuild:
ifeq ($(INSIDE_DOCKER),1)
	@echo -n '> '
	ha apps rebuild local_$(TARGET) && ha apps start local_$(TARGET)
else
	docker compose exec -T devcontainer make TARGET=$(TARGET) rebuild
endif
.PHONY: rebuild

logs:
ifeq ($(INSIDE_DOCKER),1)
	journalctl --no-tail -f -u hassio-supervisor -u hassio-bootstrap -u hassio-apparmor -u docker
else
	docker compose exec devcontainer make logs
endif
.PHONY: dev
