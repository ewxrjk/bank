bindir=/usr/local/bin
INSTALL=install

all: check

bank: $(wildcard *.go) $(wildcard cmd/bank/*.go) $(wildcard */*.go) cmd/bank/ui.go
	go build -o $@ ./cmd/bank

embed: $(wildcard cmd/embed/*.go)
	go build -o $@ ./cmd/embed

check: bank
	go test -v ./...
	./gbtest.py

EMBED=$(sort $(wildcard ui/*.html ui/*.png) ui/app.js ui/app.css)
cmd/bank/ui.go: ${EMBED} Makefile embed
	./embed -o $@ -p main ${EMBED}

install:
	$(INSTALL) -m 555 bank $(bindir)/bank

clean:
	rm -f bank embed
