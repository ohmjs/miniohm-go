all: build

.PHONY: build
build: build/testmain

.PHONY: test
test: build/testmain testdata/arithmetic.wasm
	build/testmain -wasm testdata/arithmetic.wasm -input "1+2*3-4/2"

testdata/arithmetic.wasm: testdata/arithmetic.ohm
	npx --package @ohm-js/wasm ohm2wasm testdata/arithmetic.ohm

build/testmain: testmain.go matcher.go
	go mod tidy && go build -o build/testmain

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
