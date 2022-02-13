GOPATH:=$(shell go env GOPATH)
bindir=/usr/local/bin
INSTALL=install
EMBED=$(sort $(wildcard cmd/bank/web/*.html cmd/bank/web/*.png) cmd/bank/web/app.js cmd/bank/web/app.css)

all: check

bank: $(wildcard *.go) $(wildcard cmd/bank/*.go) $(wildcard */*.go) $(EMBED) version.go
	go build -o $@ ./cmd/bank

check: bank
	go test -v ./...
	scripts/gbtest.py

version.go: scripts/make-version $(wildcard *.go) $(wildcard cmd/bank/*.go) $(wildcard */*.go)
	scripts/make-version > version.go

install:
	$(INSTALL) -m 555 bank $(bindir)/bank

$(DEP):
	go get -u github.com/golang/dep/cmd/dep

clean:
	rm -f bank
