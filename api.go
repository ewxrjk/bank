package bank

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Bank wraps access to a bank database.
type Bank struct {
	DB *sql.DB
}

// ErrUserExists is returned when creating a user that already exists.
var ErrUserExists = errors.New("user exists")

// ErrAccountExists is returned when creating an account that already exists.
var ErrAccountExists = errors.New("account exists")

// ErrInsufficientFunds is returned when there are insufficient funds for some action.
var ErrInsufficientFunds = errors.New("insufficient funds")

// ErrUnsuitableParties is returned when there is something wrong with the proposed parties to a transaction.
var ErrUnsuitableParties = errors.New("invalid or inconsistent parties to transaction")

// NewBank creates a new Bank object.
func NewBank(driver, source string) (b *Bank, err error) {
	b = &Bank{}
	if b.DB, err = OpenDatabase(driver, source); err != nil {
		err = fmt.Errorf("opening database: %v", err)
		return
	}
	return
}

// Close closes the Bank object.
func (b *Bank) Close() (err error) {
	if err = b.DB.Close(); err != nil {
		err = fmt.Errorf("closing database: %v", err)
		return
	}
	return
}

// Table creation

// NewBank creates tables for a new bank.
func (b *Bank) NewBank() (err error) {
	err = Transact(b.DB, func(tx *sql.Tx) (err error) {
		if _, err = tx.Exec(`CREATE TABLE users (
			user TEXT PRIMARY KEY,
			password TEXT
		  )`); err != nil {
			return
		}
		if _, err = tx.Exec(`CREATE TABLE transactions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			time INTEGER,
			user TEXT,
			description TEXT,
			origin TEXT,
			destination TEXT,
			amount INTEGER,
			origin_balance_after INTEGER,
			destination_balance_after INTEGER
		  )`); err != nil {
			return
		}
		if _, err = tx.Exec(`CREATE TABLE accounts (
			user TEXT PRIMARY KEY,
			balance INTEGER
		  )`); err != nil {
			return
		}
		if _, err = tx.Exec(`CREATE TABLE config (
			key TEXT PRIMARY KEY,
			value TEXT
		  )`); err != nil {
			return
		}
		return
	})
	if err != nil {
		err = fmt.Errorf("creating tables: %v", err)
	}
	return
}

// User management

// NewUser creates a new user.
func (b *Bank) NewUser(user, password string) (err error) {
	u := User{user, ""}
	if u.Password, err = SetPassword(password); err != nil {
		err = fmt.Errorf("setting new user password: %v", err)
		return
	}
	err = Transact(b.DB, func(tx *sql.Tx) (err error) {
		if err = (&User{}).Get(tx, user); err == nil {
			err = ErrUserExists
			return
		}
		return u.Put(tx, true)
	})
	if err != nil && err != ErrUserExists {
		err = fmt.Errorf("setting new user password: %v", err)
	}
	return
}

// GetUsers list users.
func (b *Bank) GetUsers() (users []string, err error) {
	err = Transact(b.DB, func(tx *sql.Tx) (err error) {
		if users, err = GetUsers(tx); err != nil {
			return
		}
		return
	})
	if err != nil {
		err = fmt.Errorf("getting users: %v", err)
	}
	return
}

// DeleteUser deletes a user
func (b *Bank) DeleteUser(user string) (err error) {
	err = Transact(b.DB, func(tx *sql.Tx) (err error) {
		if err = (&User{User: user}).Delete(tx); err != nil {
			return
		}
		return
	})
	if err != nil && err != ErrNoSuchUser {
		err = fmt.Errorf("deleting user: %v", err)
	}
	return
}

// SetPassword changes a user password.
func (b *Bank) SetPassword(user, password string) (err error) {
	err = Transact(b.DB, func(tx *sql.Tx) (err error) {
		u := User{}
		if err = u.Get(tx, user); err != nil {
			return
		}
		if u.Password, err = SetPassword(password); err != nil {
			return
		}
		return u.Put(tx, false)
	})
	if err != nil && err != ErrNoSuchUser {
		err = fmt.Errorf("changing user password: %v", err)
	}
	return
}

// CheckPassword checks a user password.
func (b *Bank) CheckPassword(user, password string) (err error) {
	err = Transact(b.DB, func(tx *sql.Tx) (err error) {
		u := User{}
		if err = u.Get(tx, user); err != nil {
			return
		}
		if err = VerifyPassword(u.Password, password); err != nil {
			return
		}
		return
	})
	if err != nil && err != ErrPasswordMismatch && err != ErrNoSuchUser {
		err = fmt.Errorf("verifying password: %v", err)
	}
	return
}

// Account management

// NewAccount creates a new account.
func (b *Bank) NewAccount(account string) (err error) {
	err = Transact(b.DB, func(tx *sql.Tx) (err error) {
		if err = (&Account{Account: account}).Get(tx, account); err == nil {
			err = ErrAccountExists
			return
		}
		a := Account{account, 0}
		return a.Put(tx, true)
	})
	if err != nil && err != ErrAccountExists {
		err = fmt.Errorf("creating account: %v", err)
	}
	return
}

// GetAccounts list accounts.
func (b *Bank) GetAccounts() (accounts []string, err error) {
	err = Transact(b.DB, func(tx *sql.Tx) (err error) {
		if accounts, err = GetAccounts(tx); err != nil {
			return
		}
		return
	})
	if err != nil {
		err = fmt.Errorf("getting accounts: %v", err)
	}
	return
}

// DeleteAccount deletes a user
func (b *Bank) DeleteAccount(account string) (err error) {
	err = Transact(b.DB, func(tx *sql.Tx) (err error) {
		if err = (&Account{Account: account}).Delete(tx); err != nil {
			return
		}
		return
	})
	if err != nil && err != ErrAccountHasBalance && err != ErrNoSuchAccount {
		err = fmt.Errorf("deleting account: %v", err)
	}
	return
}

// Transaction management

// GetTransactions gets some transactions.
//
// Starting with the most recent transaction, it skips `offset` transactions
// and then gets the next `limit` transactions.
func (b *Bank) GetTransactions(limit, offset int, after int) (transactions []Transaction, err error) {
	err = Transact(b.DB, func(tx *sql.Tx) (err error) {
		if transactions, err = GetTransactions(tx, limit, offset, after); err != nil {
			return
		}
		return
	})
	if err != nil {
		err = fmt.Errorf("getting transactions: %v", err)
	}
	return
}

// NewTransaction makes a transaction.
func (b *Bank) NewTransaction(user, origin, destination, description string, amount int) (err error) {
	err = Transact(b.DB, func(tx *sql.Tx) (err error) {
		var o, d Account
		if o.Get(tx, origin); err != nil {
			return
		}
		if d.Get(tx, destination); err != nil {
			return
		}
		if err = newTransaction(tx, user, &o, &d, description, amount); err != nil {
			return
		}
		return
	})
	if err != nil && err != ErrUnsuitableParties && err != ErrNoSuchAccount {
		err = fmt.Errorf("creating transaction: %v", err)
	}
	return
}

// newTransaction make a raw transaction
func newTransaction(tx *sql.Tx, user string, o, d *Account, description string, amount int) (err error) {
	if o.Account == d.Account {
		err = ErrUnsuitableParties
		return
	}
	if amount < 0 {
		return newTransaction(tx, user, d, o, description, -amount)
	}
	o.Balance -= amount
	if err = o.Put(tx, false); err != nil {
		return
	}
	d.Balance += amount
	if err = d.Put(tx, false); err != nil {
		return
	}
	t := Transaction{
		Time:                    time.Now().Unix(),
		User:                    user,
		Origin:                  o.Account,
		Destination:             d.Account,
		Description:             description,
		Amount:                  amount,
		OriginBalanceAfter:      o.Balance,
		DestinationBalanceAfter: d.Balance,
	}
	if err = t.Put(tx, true); err != nil {
		return
	}
	return
}

// Distribute makes a distribution transaction.
func (b *Bank) Distribute(user string, origin string, destinations []string, description string) (err error) {
	if len(destinations) < 2 {
		err = ErrUnsuitableParties
		return
	}
	for _, destination := range destinations {
		if origin == destination {
			err = ErrUnsuitableParties
			return
		}
	}
	err = Transact(b.DB, func(tx *sql.Tx) (err error) {
		var o Account
		if o.Get(tx, origin); err != nil {
			return
		}
		amount := o.Balance / len(destinations)
		if amount == 0 {
			err = ErrInsufficientFunds
			return
		}
		for _, destination := range destinations {
			var d Account
			if d.Get(tx, destination); err != nil {
				return
			}
			if err = newTransaction(tx, user, &o, &d, description, amount); err != nil {
				return
			}
		}
		return
	})
	if err != nil && err != ErrInsufficientFunds && err != ErrNoSuchAccount && err != ErrUnsuitableParties {
		err = fmt.Errorf("creating distribution transaction: %v", err)
	}
	return
}

// GetConfig retrieves a configuration item.
func (b *Bank) GetConfig(key string) (value string, err error) {
	err = Transact(b.DB, func(tx *sql.Tx) (err error) {
		value, err = GetConfig(tx, key)
		return
	})
	if err != nil && err != ErrNoSuchConfig {
		err = fmt.Errorf("getting configuration: %v", err)
	}
	return
}

// PutConfig sets a configuration item.
func (b *Bank) PutConfig(key, value string) (err error) {
	err = Transact(b.DB, func(tx *sql.Tx) (err error) {
		err = PutConfig(tx, key, value)
		return
	})
	if err != nil && err != ErrNoSuchConfig {
		err = fmt.Errorf("putting configuration: %v", err)
	}
	return
}

// GetConfigs retrieves all configuration items.
func (b *Bank) GetConfigs() (configs map[string]string, err error) {
	err = Transact(b.DB, func(tx *sql.Tx) (err error) {
		configs, err = GetConfigs(tx)
		return
	})
	if err != nil && err != ErrNoSuchConfig {
		err = fmt.Errorf("getting configuration: %v", err)
	}
	return
}
