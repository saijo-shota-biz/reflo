.PHONY: start build test tidy generate

GOFLAGS = -trimpath

start:
	go run $(GOFLAGS) ./cmd/reflo start

build:
	go build $(GOFLAGS) -o bin/reflo ./cmd/reflo

test:
	go test ./...

generate:
	go generate ./...

tidy:
	go mod tidy
