package bank

import (
	"errors"
	// sqlite3 "github.com/mattn/go-sqlite3"
	"database/sql"
)

// Account defines a single account.
type Account struct {
	// Name of the account
	Account string

	// Current balance
	Balance int
}

// ErrNoSuchAccount is returned when accessing an account whic does not exist.
var ErrNoSuchAccount = errors.New("account does not exist")

// GetAccounts returns the list of account names.
func GetAccounts(tx *sql.Tx) (accounts []string, err error) {
	var rows *sql.Rows
	if rows, err = tx.Query("SELECT user FROM accounts ORDER BY user"); err != nil {
		return
	}
	defer rows.Close()
	accounts = []string{}
	for rows.Next() {
		var account string
		if err = rows.Scan(&account); err != nil {
			return
		}
		accounts = append(accounts, account)
	}
	return
}

// Get fills in the data for a account.
func (a *Account) Get(tx *sql.Tx, account string) (err error) {
	err = tx.QueryRow("SELECT user, balance FROM accounts WHERE user=?", account).Scan(&a.Account, &a.Balance)
	if err == sql.ErrNoRows {
		err = ErrNoSuchAccount
	}
	if err != nil {
		return
	}
	return
}

// Put updates account data in the database.
// Set new to true to create or false to update.
func (a *Account) Put(tx *sql.Tx, new bool) (err error) {
	if new {
		_, err = tx.Exec("INSERT INTO accounts (user, balance) VALUES (?, ?)", a.Account, a.Balance)
	} else {
		_, err = tx.Exec("UPDATE accounts SET balance=? WHERE user=?", a.Balance, a.Account)
	}
	if err != nil {
		return
	}
	return
}
