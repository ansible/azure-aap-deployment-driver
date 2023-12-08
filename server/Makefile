# checking for go on this line
ifeq (, $(shell which go))
$(error "No go in $(PATH), can't go anywhere...")
endif
GOVERSION := $(shell go version)

.PHONY: build test test-with-coverage lint print-version

all: test build

print-version:
	@echo "Using $(GOVERSION)"

build: print-version
	mkdir -p build && go build -o build/server .

test: print-version
	go test -cover -count=1 ./...

test-with-coverage: print-version
	go test -coverpkg=./... -coverprofile=precoverage.txt -count=1 ./...
	@cat precoverage.txt | grep -v test > coverage.txt && rm precoverage.txt

lint: print-version
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOPATH)/bin/golangci-lint run -v
