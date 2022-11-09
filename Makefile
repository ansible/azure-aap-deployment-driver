export SHELL := /usr/bin/bash
# checking for go on this line
ifeq (, $(shell which go))
$(error "No go in $(PATH), can't go anywhere...")
endif
GOVERSION := $(shell go version)

.PHONY: build test print-version

all: test build

print-version:
	@echo "Using go version $(GOVERSION)"

build: print-version
	mkdir -p build && cd server && go build -o ../build/server .

test: print-version
	cd server && go test ./...
