.PHONY: build test

build:
	go build -o bin/rowxml ./cmd/rowxml

test:
	go test ./... -v

