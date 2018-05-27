package bank

import (
	"errors"
	// sqlite3 "github.com/mattn/go-sqlite3"
	"database/sql"
)

// User defines a single user.
type User struct {
	// Name of the user
	User string

	// Hashed password
	Password string
}

// ErrNoSuchUser is returned when accessing a user who does not exist.
var ErrNoSuchUser = errors.New("user does not exist")

// GetUsers returns the list of user names.
func GetUsers(tx *sql.Tx) (users []string, err error) {
	var rows *sql.Rows
	if rows, err = tx.Query("SELECT user FROM users ORDER BY user"); err != nil {
		return
	}
	defer rows.Close()
	users = []string{}
	for rows.Next() {
		var user string
		if err = rows.Scan(&user); err != nil {
			return
		}
		users = append(users, user)
	}
	return
}

// Get fills in the data for a user.
func (u *User) Get(tx *sql.Tx, user string) (err error) {
	err = tx.QueryRow("SELECT user, password FROM users WHERE user=?", user).Scan(&u.User, &u.Password)
	if err == sql.ErrNoRows {
		err = ErrNoSuchUser
	}
	if err != nil {
		return
	}
	return
}

// Put updates user data in the database.
// Set new to true to create or false to update.
func (u *User) Put(tx *sql.Tx, new bool) (err error) {
	if new {
		_, err = tx.Exec("INSERT INTO users (user, password) VALUES (?, ?)", u.User, u.Password)
	} else {
		_, err = tx.Exec("UPDATE users SET password=? WHERE user=?", u.Password, u.User)
	}
	if err != nil {
		return
	}
	return
}
