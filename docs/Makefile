.PHONY: build-wasm
build-wasm:
	cp $(shell go env GOROOT)/misc/wasm/wasm_exec.js  .
	GOOS=js GOARCH=wasm go build -o main.wasm main.go
