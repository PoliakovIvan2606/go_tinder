.PHONY: run
run:
	go run ./cmd/apiserver

.PHONY: test
test:test
	go test -v -race -timeout 30s ./...
.PHONY: build
build:build
	go build ./cmd/apiserver
	
.DEFAULT_GOAL := run