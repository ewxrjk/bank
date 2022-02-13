package bank

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3" // include sqlite3 driver
	"log"
)

// Transact runs a function inside a database transaction.
// On success the transaction is committed.
// On error the transaction is rolled back.
func Transact(database *sql.DB, body func(tx *sql.Tx) (err error)) (err error) {
	var tx *sql.Tx
	if tx, err = database.Begin(); err != nil {
		return
	}
	if err = body(tx); err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			log.Fatalf("rollback: %v", err2)
		}
		return
	}
	return tx.Commit()
}

// OpenDatabase opens a database handle.
func OpenDatabase(driver, source string) (db *sql.DB, err error) {
	if db, err = sql.Open(driver, source); err != nil {
		return
	}
	if driver == "sqlite3" {
		if _, err = db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
			log.Printf("enabling foreign keys: %s", err)
			if err = db.Close(); err != nil {
				log.Printf("closing database: %s", err)
			}
			db = nil
			return
		}
	} else {
		log.Printf("warning: unrecognized database driver: %s", driver)
	}
	return
}
