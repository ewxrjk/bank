# Bank

This is web application for tracking shared finances within a small, mutually trusting group.

# Installation

TODO

# Setup

## Database Setup

Bank uses [Sqlite](https://www.sqlite.org/) as its database,
meaning that all data is stored in a single file.
This file must have read-write permission for the server.

To create an empty database,
supposing that it is called `bank.db`
and located in the current directory:

    $ gobank -d bank.db init

You will need to create an initial user:

    $ gobank -d bank.db user add rjk 
    Enter password: 
    Confirm password: 

You can create additional users at this stage,
or via the web interface.

Passwords may subsequently be changed either with `gobank user pw` or via the web interface.

## Server Setup

TODO

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
Values are substitued into HTML pages using [html/template](https://golang.org/pkg/html/template/)
and into database queries using placeholder parameters as described in the [database/sql](https://golang.org/pkg/database/sql/) API.

User logins are tracked with a cookie and an associated token.
The token is embedded into HTML pages
and must be presented with all mutating operations.
Non-mutating requests require only the cookie.
