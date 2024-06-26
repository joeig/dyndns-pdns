GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOCOVER=$(GOCMD) tool cover
GOGET=$(GOCMD) get
GOFMT=gofmt
BINARY_NAME=dyndns-pdns
GOFILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*")

.DEFAULT_GOAL := all
.PHONY: all build build-linux-amd64 test coverage check-fmt fmt clean run

all: check-fmt test coverage build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/dyndns-pdns

build-linux-amd64:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)_linux_amd64 -v ./cmd/dyndns-pdns

test:
	$(GOTEST) -v ./... -covermode=count -coverprofile=coverage.out

coverage:
	$(GOCOVER) -func=coverage.out

check-fmt:
	$(GOFMT) -d ${GOFILES}

fmt:
	$(GOFMT) -w ${GOFILES}

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)_linux_amd64

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)
