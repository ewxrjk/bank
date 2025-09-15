bindir=/usr/local/bin
INSTALL=install
EMBED=$(sort $(wildcard cmd/bank/web/*.html cmd/bank/web/*.png) cmd/bank/web/app.js cmd/bank/web/app.css)

all: check

bank: $(wildcard *.go) $(wildcard cmd/bank/*.go) $(wildcard */*.go) $(EMBED) pkg/bank/version.go
	go build -o $@ ./cmd/bank

check: bank
	go test -v ./...
	scripts/gbtest.py

pkg/bank/version.go: scripts/make-version $(wildcard *.go) $(wildcard cmd/bank/*.go) $(wildcard */*.go)
	scripts/make-version > pkg/bank/version.go

install:
	$(INSTALL) -m 555 bank $(bindir)/bank

clean:
	rm -f bank
