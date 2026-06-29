BINARY  := paperless
PREFIX  ?= /usr/local
DESTDIR ?=

.PHONY: build install generate clean

build:
	go build -buildvcs=false -o $(BINARY) .

install: build
	install -d $(DESTDIR)$(PREFIX)/bin
	install -m 755 $(BINARY) $(DESTDIR)$(PREFIX)/bin/$(BINARY)
	@echo "Installed to $(DESTDIR)$(PREFIX)/bin/$(BINARY)"

generate:
	@echo "Downloading schema..."
	curl -sf "$$PAPERLESS_URL/api/schema/?format=json" \
		-H "Authorization: Token $$PAPERLESS_API_TOKEN" \
		-o schema/paperless.json
	@echo "Fixing schema inconsistencies..."
	python3 scripts/fix-schema.py
	@echo "Generating API client..."
	oapi-codegen --config oapi-codegen.yaml schema/paperless.json
	@echo "Updating tracked version..."
	curl -sf https://api.github.com/repos/paperless-ngx/paperless-ngx/releases/latest \
		| python3 -c "import json,sys; print(json.load(sys.stdin)['tag_name'])" \
		> .paperless-version
	@echo "Done. Tracked version: $$(cat .paperless-version)"

clean:
	rm -f $(BINARY)
