all: build

.PHONY: build
build: build/example

.PHONY: test
test: build/example testdata/arithmetic.wasm
	build/example -wasm testdata/arithmetic.wasm -input "1+2*3-4/2" -verbose

testdata/arithmetic.wasm: testdata/arithmetic.ohm
	npx --package @ohm-js/wasm ohm2wasm testdata/arithmetic.ohm

build/example: cmd/example/main.go matcher.go cst.go
	go mod tidy && go build -o build/example ./cmd/example

.PHONY: tag
tag:
	@echo "Most recent tag: $$(git describe --tags --abbrev=0 2>/dev/null || echo 'No tags yet')"
	@if git diff-index --quiet HEAD --; then \
	  echo "Creating new tag..."; \
		read -p "Enter tag name (e.g., v1.0.0): " tag_name; \
		git tag $$tag_name; \
		git push origin $$tag_name; \
		echo "Tag $$tag_name created and pushed."; \
	else \
	  echo "Error: Working directory not clean. Commit changes before tagging."; \
		exit 1; \
	fi
