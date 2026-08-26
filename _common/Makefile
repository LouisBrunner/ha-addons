PORT ?= 8787
DIST_DIR = dist
TARBALL = $(DIST_DIR)/common.tar.gz

serve: $(TARBALL)
	python3 -m http.server $(PORT) --directory $(DIST_DIR)
.PHONY: serve

$(TARBALL): $(DIST_DIR) $(shell find rootfs -type f)
	COPYFILE_DISABLE=1 tar czf $@ -C rootfs .

$(DIST_DIR):
	mkdir -p $@

clean:
	rm -rf $(DIST_DIR)
.PHONY: clean
