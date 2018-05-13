bankdir=/var/lib/bank
wwwdir=/var/www/bank
testwwwdir=/var/www/testbank
INSTALL=install

all: bank.real bank.site.real

gobank: $(wildcard *.go) $(wildcard cmd/bank/*.go) cmd/bank/ui.go
	go build -o $@ ./cmd/bank

embed: $(wildcard cmd/embed/*.go)
	go build -o $@ ./cmd/embed

gocheck: gobank
	go test -v ./...
	./gbtest.py

EMBED=$(wildcard ui/*.html) ui/app.js ui/app.css
cmd/bank/ui.go: ${EMBED} Makefile embed
	./embed -o $@ -p main ${EMBED}

bank.real: bank
	rm -f bank.real
	sed < bank > bank.real s/testbank/bank/g;
	chmod 555 bank.real

bank.site.real: bank.site
	rm -f bank.site.real
	sed < bank.site > bank.site.real s/testbank/bank/g;
	chmod 444 bank.site.real

check:
	perl -wc bank

install-real: check
	adduser --system --group --home $(bankdir) bank
	chmod 700 $(bankdir)
	mkdir -m 755 -p $(wwwdir)
	chown bank:bank $(wwwdir)
	$(INSTALL) favicon.ico $(wwwdir)/favicon.ico
	$(INSTALL) -o bank -g bank -m 755 bank.real $(wwwdir)/bank
	$(INSTALL) bank.site.real /etc/apache2/sites-available/bank
	ln -sf ../sites-available/bank /etc/apache2/sites-enabled
	mkdir -m 755 -p /var/log/apache2/bank
	service apache2 reload

install-test:
	mkdir -m 755 -p $(testwwwdir)
	chown bank:bank $(testwwwdir)
	$(INSTALL) favicon.ico $(testwwwdir)/favicon.ico
	$(INSTALL) -o bank -g bank -m 755 bank $(testwwwdir)/bank
	$(INSTALL) bank.site /etc/apache2/sites-available/testbank
	ln -sf ../sites-available/testbank /etc/apache2/sites-enabled
	mkdir -m 755 -p /var/log/apache2/testbank
	service apache2 reload

setup-real: check
	su bank -s $(SHELL) -c "sqlite3 -init bank.sql $(bankdir)/bank.db < /dev/null"
	chmod 600 $(bankdir)/bank.db

setup-test: check
	su bank -s $(SHELL) -c "sqlite3 -init bank.sql $(bankdir)/testbank.db < /dev/null"
	chmod 600 $(bankdir)/testbank.db

clean:
	rm -f bank.real
	rm -f bank.site.real

# TODO logfile rotation
