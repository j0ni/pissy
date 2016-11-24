GOFILES = $(shell find . -name '*.go' -not -path './vendor/*')
GOPACKAGES = $(shell go list ./...  | grep -v /vendor/)

default: build check

test:
	@echo go test -v -race $(TEST_TAGS) ./...
	@go test -v -race $(TEST_TAGS) $(GOPACKAGES)

vet:
	@echo go vet
	@go vet $(GOPACKAGES)

fmt:
	@echo gofmt
	@if gofmt -l $(GOFILES) | grep .; then echo "Code differs from gofmt's style" 1>&2 && exit 1; fi

build:
	@echo go build
	@go build

clean:
	@rm -f pissy
	@rm -rf vendor/*

check: test vet fmt
