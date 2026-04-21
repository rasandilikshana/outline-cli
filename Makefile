BINARY_NAME=outline
MCP_BINARY_NAME=outline-mcp
VERSION=0.3.0
INSTALL_DIR=/usr/local/bin

LDFLAGS_CLI=-ldflags="-X outline-cli/cmd.Version=$(VERSION)"
LDFLAGS_MCP=-ldflags="-X main.Version=$(VERSION)"

.PHONY: build build-mcp build-all clean install install-mcp install-all run fmt vet tidy test

## Build targets

build:
	go build $(LDFLAGS_CLI) -o $(BINARY_NAME) ./main.go

build-mcp:
	go build $(LDFLAGS_MCP) -o $(MCP_BINARY_NAME) ./cmd/outline-mcp

build-all: build build-mcp

## Install targets

install: build
	cp $(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)

install-mcp: build-mcp
	cp $(MCP_BINARY_NAME) $(INSTALL_DIR)/$(MCP_BINARY_NAME)

install-all: install install-mcp

## Housekeeping

clean:
	rm -f $(BINARY_NAME) $(MCP_BINARY_NAME)

run: build
	./$(BINARY_NAME)

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy

test:
	go test ./...
