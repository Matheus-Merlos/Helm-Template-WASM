.PHONY: build clean copy-runtime

all: build

build: dist
	CGO_ENABLED=0 GOOS=js GOARCH=wasm go build \
		-tags osusergo \
		-ldflags="-s -w" \
		-o dist/helm-template.wasm .

	@echo "Built dist/helm-template.wasm"

copy-runtime: dist
	cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" dist/
	@echo "Copied wasm_exec.js to dist/"

dist:
	mkdir -p dist

clean:
	rm -rf dist

release: build copy-runtime