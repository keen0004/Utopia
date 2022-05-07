# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: build test clean all

# os can be linux/darwin/freebsd/windows
# arch can be 386/amd64/arm
OS=darwin
ARCH=amd64

GO=go
ENV=env GO111MODULE=on GOOS=$(OS) GOARCH=$(ARCH)
GOBUILD=$(ENV) $(GO) build
GORUN=$(ENV) $(GO) run
FLAGS=-v
BINDIR=./build/bin/

build:
	$(GOBUILD) $(FLAGS) -o $(BINDIR)/chaintool ./cmd/chaintool/
	$(GOBUILD) $(FLAGS) -o $(BINDIR)/cointool ./cmd/cointool/
	$(GOBUILD) $(FLAGS) -o $(BINDIR)/contracttool ./cmd/contracttool/
	$(GOBUILD) $(FLAGS) -o $(BINDIR)/keytool ./cmd/keytool/
	$(GOBUILD) $(FLAGS) -o $(BINDIR)/offertool ./cmd/offertool/

all: build test

test: 
	$(go) test ./test

clean:
	$(GO) clean -cache
	rm -fr build/bin/*

