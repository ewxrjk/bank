GOPATH:=$(shell go env GOPATH)
bindir=/usr/local/bin
INSTALL=install
DEP=$(GOPATH)/bin/dep

all: check

bank: $(wildcard *.go) $(wildcard cmd/bank/*.go) $(wildcard */*.go) cmd/bank/ui.go version.go vendor
	go build -o $@ ./cmd/bank

embed: $(wildcard cmd/embed/*.go) vendor
	go build -o $@ ./cmd/embed

check: bank
	go test -v ./...
	./gbtest.py

vendor: $(DEP)
	$(DEP) ensure

EMBED=$(sort $(wildcard ui/*.html ui/*.png) ui/app.js ui/app.css)
cmd/bank/ui.go: ${EMBED} Makefile embed
	./embed -o $@ -p main ${EMBED}

version.go: scripts/make-version $(wildcard *.go) $(wildcard cmd/bank/*.go) $(wildcard */*.go) cmd/bank/ui.go
	scripts/make-version > version.go

install:
	$(INSTALL) -m 555 bank $(bindir)/bank

$(DEP):
	go get -u github.com/golang/dep/cmd/dep

clean:
	rm -f bank embed
