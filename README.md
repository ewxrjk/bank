# Description

This is web application for tracking shared finances within a small, mutually trusting group.

# Setup

## Installation

You will need:

- [Go](https://golang.org/)
- [dep](https://github.com/golang/dep)
- [Sqlite](https://www.sqlite.org/)

To build and self-test the software:

    $ go get github.com/ewxrjk/bank
    $ cd ${GOPATH:-$HOME/go}/src/github.com/ewxrjk/bank
    $ dep ensure
    $ make

To install the `bank` command in `/usr/local/bin`:

    $ sudo make install

## Database Setup

Bank uses [Sqlite](https://www.sqlite.org/) as its database,
meaning that all data is stored in a single file.
This file must have read-write permission for the server.

To create an empty database,
supposing that it is called `bank.db`
and located in the current directory:

    $ bank -d bank.db init

You will need to create an initial user:

    $ bank -d bank.db user add rjk
    Enter password:
    Confirm password:

You can create additional users at this stage,
or via the web interface.

Passwords may subsequently be changed either with `bank user pw` or via the web interface.

You can set the name of the bank.
This is used in the web interface.

    $ bank -d bank.db config set title 'Example Bank'

## Server Setup

To run a server without TLS:

    $ bank -d bank.db server --address :8080

### Caching

You can use the `--lifetime` option to set the maximum caching period for static content
(which currently means the CSS and JavaScript).
The default is 60 seconds.
For a production system, where these files may be unchanged for weeks or months on end,
a larger value may be suitable.

### TLS

You can use the `--cert` and `--key` options to specify a TLS certificate and key.
[Let's Encrypt](https://letsencrypt.org/) is the best way to acquire a certificate,
but you can get started with a self-signed certificate as follows:

    $ openssl req -new -newkey rsa:2048 -x509 -nodes -subj /CN=www.example.com -keyout key.pem -out cert.pem
    $ bank -d bank.db server --address :8080 --cert cert.pem --key key.pem

## Deployment

I run the service as its own user/group:

    # useradd -Urmd/var/lib/bank bank
    # su -lcbash bank
    $ bank init
    $ bank user add rjk
    Enter password:
    Confirm password:

Administrative users may be added to the `bank` group.

To run it as a daemon,
edit `bank.service` and install and enable it:

    # install -m644 bank.service /usr/local/lib/systemd/system/bank.service
    # systemctl enable --now bank.service
    Created symlink /etc/systemd/system/multi-user.target.wants/bank.service → /usr/local/lib/systemd/system/bank.service.

I front-end the service with Apache,
so I can take care of TLS in a uniform way
with other services.
The following directives forward requests to the service:

	ProxyPass "/" "http://localhost:8344/"
	ProxyPassReverse "/" "http://localhost:8344/"

`mod_proxy` and `mod_proxy_http` must be enabled.

# Usage

A typical setup would have:

- an account for each user, reflecting how much they are owed
- an account called `house` reflecting how much the collective of users are owed.
Normally `house` would have a negative balance, reflecting that the collective owes money to its individual members.

The common actions in this model are:

- payments from `house` to some individual's account, reflecting that the individual has made a payment on behalf of the collective.
This can be done via the _New Transaction_ page, or via the form on the front page.
- distribution of the `house` balance among the other accounts,
reflecting that the collective's obligations full upon to its members.
This can be done via the _Distribute_ page.
- payments between individual users, reflecting offline resolution of these obligations.
This can be done via the _New Transaction_ page, or via the form on the front page.

## Example

Suppose the users are `fred` and `bob`, and they are populating their kitchen.

1. `fred` buys a toaster for £20.
He fills in the details of this real-life transaction on the front page,
causing a payment from `house` to `fred` of £20 to be recorded.
2. `bob` buys a microwave oven for £40.
He fills in the details on the front page,
causing a payment from `house` to `bob` of £60 to be recorded.
3. `fred` now has a balance of £20 and `bob` of £40.
`house` has a balance of -£60.
4. Any user uses the _Distribute_ page
to distribute from `house` to `fred` and `bob`.
This divides its balance into two and transfers each half to one of the human users.
The effect is that `fred` has a balanced of -£10 and `bob` of £10.
5. `fred` owes £10, and `bob` is owed £10, so (in real life) Fred hands Bob £10.
6. Either user enters the details of this real-life transaction via the front page,
causing a payment of £10 from `bob` to `fred` to be recorded,
leaving each with a balance of £0.

# Security

All legitimate users are treated equally
and can change one another's
passwords, enter transactions on behalf of
each other, etc.

Passwords are obscured with [scrypt](https://en.wikipedia.org/wiki/Scrypt).
The parameters are stored in the database to make changing them easy.
Please see the source code for the current default parameters.

The server is implemented in [Go](https://golang.org/),
a memory-safe language.
Values are substituted into HTML pages using [html/template](https://golang.org/pkg/html/template/)
and into database queries using placeholder parameters as described in the [database/sql](https://golang.org/pkg/database/sql/) API.

User logins are tracked with a cookie and an associated token.
The token is embedded into HTML pages
and must be presented with all mutating operations.
Non-mutating requests require only the cookie.

The JavaScript UI uses [jQuery](https://jquery.com/) heavily.
It constrains the jQuery code using sub-resource integrity.
