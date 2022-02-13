package bank

import (
	"database/sql"
)

// Transaction defines a single transaction.
type Transaction struct {
	// Transaction ID
	ID int

	// Unix time of the transaction
	Time int64

	// User who make the transaction
	User string

	// Description of the transaction
	Description string

	// Origin account
	Origin string

	// Destination account
	Destination string

	// Amount in pence
	Amount int

	OriginBalanceAfter int

	DestinationBalanceAfter int
}

// GetTransactions returns all transactions in a range.
func GetTransactions(tx *sql.Tx, limit, offset int, after int) (transactions []Transaction, err error) {
	var rows *sql.Rows
	query := "SELECT id,time,user,description,origin,destination,amount,origin_balance_after,destination_balance_after FROM transactions WHERE id > ? ORDER BY id DESC"
	args := []interface{}{after}
	if limit != 0 {
		query = query + " LIMIT ?"
		args = append(args, limit)
	}
	if offset != 0 {
		query = query + " OFFSET ?"
		args = append(args, offset)
	}
	if rows, err = tx.Query(query, args...); err != nil {
		return
	}
	defer rows.Close()
	transactions = []Transaction{}
	for rows.Next() {
		var transaction Transaction
		if err = rows.Scan(&transaction.ID, &transaction.Time, &transaction.User, &transaction.Description, &transaction.Origin, &transaction.Destination, &transaction.Amount, &transaction.OriginBalanceAfter, &transaction.DestinationBalanceAfter); err != nil {
			return
		}
		transactions = append(transactions, transaction)
	}
	return
}

// Put updates transaction data in the database.
// Set new to true to create or false to update.
func (t *Transaction) Put(tx *sql.Tx, new bool) (err error) {
	if new {
		_, err = tx.Exec("INSERT INTO transactions (time,user,description,origin,destination,amount,origin_balance_after,destination_balance_after) VALUES (?,?,?,?,?,?,?,?)",
			t.Time, t.User, t.Description, t.Origin, t.Destination, t.Amount, t.OriginBalanceAfter, t.DestinationBalanceAfter)
	} else {
		panic("NYI")
	}
	return
}
