TARGET ?= unknown

INSIDE_DOCKER = $(shell stat /.indocker 2>&1 >/dev/null && echo 1 || echo 0)
ifeq ($(INSIDE_DOCKER),1)
endif

CONFIG = $(TARGET)/config.yaml
IMAGE_COMMENT_CMD = sed -i '' s/^image:/\#image:/ $(CONFIG)
IMAGE_UNCOMMENT_CMD = sed -i '' s/^\#image:/image:/ $(CONFIG)

all:
.PHONY: all

setup:
ifeq ($(INSIDE_DOCKER),1)
	ha apps install local_$(TARGET)
else
	lifecycle -s $(IMAGE_COMMENT_CMD) % -e $(IMAGE_UNCOMMENT_CMD) % -- docker compose exec -T devcontainer make TARGET=$(TARGET) setup
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
	hot -d $(TARGET) -d _common -e dist make TARGET=$(TARGET) rebuild
endif
.PHONY: dev

dev-serve-common:
ifeq ($(INSIDE_DOCKER),1)
	@docker inspect common-serve >/dev/null 2>&1 \
		|| docker run -d --name common-serve -p 8787:80 -v $(CURDIR)/_common/dist:/usr/share/nginx/html:ro nginx:alpine >/dev/null
else
	@echo "# Unsupported outside the container"
endif
.PHONY: dev-serve-common

rebuild:
ifeq ($(INSIDE_DOCKER),1)
	@$(MAKE) dev-serve-common
	@$(MAKE) TARGET=$(TARGET) rebuild-actual
else
	@$(MAKE) -C _common dist/common.tar.gz >/dev/null
	lifecycle -s $(IMAGE_COMMENT_CMD) % -e $(IMAGE_UNCOMMENT_CMD) % -- \
		docker compose exec -T -e HOT_CHANGED_FILES devcontainer make TARGET=$(TARGET) rebuild
endif
.PHONY: rebuild

rebuild-actual:
	@case "$$HOT_CHANGED_FILES" in \
		*config.yaml*) bash _common/reset-cache.sh $(TARGET) ;; \
		*) ha apps rebuild --force local_$(TARGET) && ha apps start local_$(TARGET) ;; \
	esac
.PHONY: rebuild-actual

logs:
ifeq ($(INSIDE_DOCKER),1)
	journalctl --no-tail -f -u hassio-supervisor -u hassio-bootstrap -u hassio-apparmor -u docker
else
	docker compose exec devcontainer make logs
endif
.PHONY: dev
