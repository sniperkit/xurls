.PHONY: all test clean man fast release install

GO15VENDOREXPERIMENT=1

PROG_NAME := "xurls"

all: deps deps-ci test version build install examples

build: deps build-generate generate
	@go build -ldflags "-X main.VERSION=`cat VERSION`" -o ./bin/$(PROG_NAME) ./cmd/$(PROG_NAME)/*.go

build-generate:
	@go build -ldflags "-X main.VERSION=`cat VERSION`" -o ./bin/$(PROG_NAME)-schemesgen ./cmd/generate/schemesgen/*.go
	@go build -ldflags "-X main.VERSION=`cat VERSION`" -o ./bin/$(PROG_NAME)-tldsgen ./cmd/generate/tldsgen/*.go

generate:
	@./bin/$(PROG_NAME)-schemesgen
	@./bin/$(PROG_NAME)-tldsgen

version: deps
	@which $(PROG_NAME)
	@$(PROG_NAME) --version

install: deps build-generate generate
	@go install -ldflags "-X main.VERSION=`cat VERSION`" ./cmd/$(PROG_NAME)
	@$(PROG_NAME) --version

fast: deps
	@go build -i -ldflags "-X main.VERSION=`cat VERSION`-dev" -o ./bin/$(PROG_NAME) ./cmd/$(PROG_NAME)/*.go
	@$(PROG_NAME) --version

deps:
	@glide install --strip-vendor

deps-pkg: deps

deps-reset: deps-rm deps-init deps-pkg deps-ci

deps-rm:
	@rm -f glide.*
	@rm -fR ./vendor

deps-init:
	@yes no | glide create

deps-update:
	@glide update

deps-ci:
	@go get -v -u github.com/go-playground/overalls
	@go get -v -u github.com/mattn/goveralls
	@go get -v -u golang.org/x/tools/cmd/cover

test:
	@go test ./pkg/...

clean:
	@go clean
	@rm -fr ./bin
	@rm -fr ./dist

release: $(PROG_NAME)
	@git tag -a `cat VERSION`
	@git push origin `cat VERSION`

