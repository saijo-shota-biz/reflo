.PHONY: start build test tidy

GOFLAGS = -trimpath

start:
	go run $(GOFLAGS) ./cmd/reflo start

build:
	go build $(GOFLAGS) -o bin/reflo ./cmd/reflo

test:
	go test ./...

tidy:
	go mod tidy
