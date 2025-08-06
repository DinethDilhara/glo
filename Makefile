.PHONY: build test clean install release-local release-snapshot help

# Build variables
BINARY_NAME=glo
VERSION?=$(shell git describe --tags --always --dirty)
COMMIT?=$(shell git rev-parse --short HEAD)
DATE?=$(shell date -u '+%Y-%m-%d %H:%M:%S UTC')

# Go variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# LDFLAGS
LDFLAGS=-ldflags "-s -w -X github.com/DinethDilhara/glo/cmd.version=$(VERSION) -X github.com/DinethDilhara/glo/cmd.commit=$(COMMIT) -X 'github.com/DinethDilhara/glo/cmd.date=$(DATE)'"

## help: Display this help message
help:
	@echo "Available commands:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

## build: Build the binary
build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) .

## test: Run tests
test:
	$(GOTEST) -v ./...

## clean: Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

## install: Install the binary to /usr/local/bin
install: build
	sudo mv $(BINARY_NAME) /usr/local/bin/

## release-local: Build release locally (without publishing)
release-local:
	goreleaser release --clean --skip=publish

## release-snapshot: Build snapshot release
release-snapshot:
	goreleaser release --snapshot --clean

## deps: Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

## lint: Run golangci-lint
lint:
	golangci-lint run

## check: Run tests and linting
check: test lint

## version: Show current version
version:
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Date: $(DATE)"
