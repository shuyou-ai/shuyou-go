.PHONY: run build test tidy fmt

run:
	go run ./cmd/server -config configs/config.yaml

build:
	go build -o bin/server ./cmd/server

test:
	go test ./...

tidy:
	go mod tidy

fmt:
	go fmt ./...

setup:
	cp -n configs/config.example.yaml configs/config.yaml || true
