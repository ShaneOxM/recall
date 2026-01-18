.PHONY: build install test clean

BINARY=rc
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X github.com/shaneoxm/recall/cmd/rc/cmd.Version=$(VERSION)"

build:
	go build $(LDFLAGS) -o $(BINARY) ./cmd/rc

install:
	go install $(LDFLAGS) ./cmd/rc

test:
	go test -v ./...

clean:
	rm -f $(BINARY)
	go clean

# Cross-compilation
build-all:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/rc-darwin-amd64 ./cmd/rc
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/rc-darwin-arm64 ./cmd/rc
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/rc-linux-amd64 ./cmd/rc
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/rc-linux-arm64 ./cmd/rc
