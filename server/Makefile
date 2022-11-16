export SHELL := /usr/bin/bash
# checking for go on this line
ifeq (, $(shell which go))
$(error "No go in $(PATH), can't go anywhere...")
endif
GOVERSION := $(shell go version)

.PHONY: build test test-with-coverage print-version

all: test build

print-version:
	@echo "Using go version $(GOVERSION)"

build: print-version
	mkdir -p build && go build -o build/server .

test: print-version
	go test -cover -count=1 ./...

test-with-coverage: print-version
	go test -coverpkg=./... -coverprofile=coverage.txt -count=1 ./...
