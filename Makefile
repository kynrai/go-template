BINARIES=$$(go list ./cmd/...)
TESTABLE=$$(go list ./...)
TARGET=app

all : test build

.PHONY: deps build test

deps:
	@go get -t ./... && go mod tidy

build: 
	@go install -v $(BINARIES)

test:
	@go test -v $(TESTABLE)

local: test build
	@app