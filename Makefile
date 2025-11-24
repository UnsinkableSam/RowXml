.PHONY: build test

build:
	go build -o bin/RowXml ./cmd/RowXml

test:
	go test ./... -v

