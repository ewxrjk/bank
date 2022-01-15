bindir=/usr/local/bin
INSTALL=install
EMBED=$(sort $(wildcard cmd/bank/ui/*.html cmd/bank/ui/*.png) cmd/bank/ui/app.js cmd/bank/ui/app.css)

all: check

bank: $(wildcard *.go) $(wildcard cmd/bank/*.go) $(wildcard */*.go) $(EMBED)
	go build -o $@ ./cmd/bank

check: bank
	go test -v ./...
	./gbtest.py

install:
	$(INSTALL) -m 555 bank $(bindir)/bank

clean:
	rm -f bank
