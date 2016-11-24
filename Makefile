GOFILES = $(shell find . -name '*.go' -not -path './vendor/*')
GOPACKAGES = $(shell go list ./...  | grep -v /vendor/)

default: clean build

test:
	@echo go test -v -race $(TEST_TAGS) ./...
	@go test -v -race $(TEST_TAGS) $(GOPACKAGES)

vet:
	@echo go vet
	@go vet $(GOPACKAGES)

fmt:
	@echo gofmt
	@if gofmt -l $(GOFILES) | grep .; then echo "Code differs from gofmt's style" 1>&2 && exit 1; fi

sync: vendor/vendor.json
	@echo govendor sync
	@govendor sync

build: sync $(GOFILES)
	@echo go build
	@go build

clean:
	@echo go clean
	@go clean

check: test vet fmt

install: build
	@echo go install
	@go install

all: clean build check
