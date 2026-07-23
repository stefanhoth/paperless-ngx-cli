BINARY  := paperless
PREFIX  ?= /usr/local
DESTDIR ?=

.PHONY: build install generate clean check format lint test

build:
	go build -buildvcs=false -o $(BINARY) .

check:
	go vet ./...

format:
	gofumpt -w .

lint: check
	test -z "$$(gofumpt -l .)" || (echo "not gofumpt-formatted, run: make format" >&2 && gofumpt -l . && exit 1)
	golangci-lint run ./...

test:
	go test -race ./...

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

# Fetch schema from upstream Docker image — no running Paperless instance required.
# VERSION defaults to the latest GitHub release; override with: make generate-docker VERSION=v2.20.15
generate-docker:
	$(eval VERSION ?= $(shell curl -sf https://api.github.com/repos/paperless-ngx/paperless-ngx/releases/latest \
		| python3 -c "import json,sys; print(json.load(sys.stdin)['tag_name'])"))
	$(eval IMAGE_TAG := $(patsubst v%,%,$(VERSION)))
	@echo "Pulling ghcr.io/paperless-ngx/paperless-ngx:$(IMAGE_TAG)..."
	docker pull ghcr.io/paperless-ngx/paperless-ngx:$(IMAGE_TAG)
	@echo "Exporting schema via manage.py spectacular..."
	docker run --rm \
		--entrypoint python3 \
		-e PAPERLESS_SECRET_KEY=changeme \
		-e PAPERLESS_DBENGINE=sqlite \
		ghcr.io/paperless-ngx/paperless-ngx:$(IMAGE_TAG) \
		/usr/src/paperless/src/manage.py spectacular \
			--file /dev/stdout --format openapi-json 2>/dev/null \
		> schema/paperless.json
	@echo "Fixing schema inconsistencies..."
	python3 scripts/fix-schema.py
	@echo "Generating API client..."
	oapi-codegen --config oapi-codegen.yaml schema/paperless.json
	@echo "$(VERSION)" > .paperless-version
	@echo "Done. Tracked version: $(VERSION)"

clean:
	rm -f $(BINARY)
