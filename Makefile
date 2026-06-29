BINARY := paperless

.PHONY: build install generate clean

build:
	go build -buildvcs=false -o $(BINARY) .

install: build
	cp $(BINARY) ../bin/paperless
	@echo "Installed to bin/paperless"

generate:
	@echo "Downloading schema..."
	curl -s "$$PAPERLESS_URL/api/schema/?format=json" \
		-H "Authorization: Token $$PAPERLESS_API_TOKEN" \
		-o schema/paperless.json
	@echo "Fixing schema inconsistencies..."
	python3 scripts/fix-schema.py
	@echo "Generating API client..."
	oapi-codegen --config oapi-codegen.yaml schema/paperless.json
	@echo "Done."

clean:
	rm -f $(BINARY)
