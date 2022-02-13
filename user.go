package bank

import (

	// sqlite3 "github.com/mattn/go-sqlite3"
	"database/sql"
	"net/http"

	"github.com/ewxrjk/bank/util"
)

// User defines a single user.
type User struct {
	// Name of the user
	User string

	// Hashed password
	Password string
}

// ErrNoSuchUser is returned when accessing a user who does not exist.
var ErrNoSuchUser = util.HTTPError{"user does not exist", http.StatusNotFound}

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
	return
}

// Delete deletes a user.
func (u *User) Delete(tx *sql.Tx) (err error) {
	var r sql.Result
	if r, err = tx.Exec("DELETE FROM users WHERE user=?", u.User); err != nil {
		return
	}
	var rows int64
	if rows, err = r.RowsAffected(); err != nil {
		return
	}
	if rows != 1 {
		err = ErrNoSuchUser
		return
	}
	return
}
