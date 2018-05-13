package bank

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestApiUser(t *testing.T) {
	var err error
	var b *Bank
	os.Remove("_user.db")
	b, err = NewBank("sqlite3", "_user.db")
	require.NoError(t, err)
	defer func() {
		b.Close()
		os.Remove("_user.db")
	}()
	require.NoError(t, b.NewBank())

	assert.NoError(t, b.NewUser("fred", "fredpass"))
	assert.NoError(t, b.CheckPassword("fred", "fredpass"))
	assert.Equal(t, ErrPasswordMismatch, b.CheckPassword("fred", "wrongpass"))
	assert.Equal(t, ErrPasswordMismatch, b.CheckPassword("fred", ""))
	assert.NotNil(t, b.CheckPassword("bob", "fredpass")) // TODO ugly

	assert.NoError(t, b.SetPassword("fred", "newpass"))
	assert.NoError(t, b.CheckPassword("fred", "newpass"))
	assert.Equal(t, ErrPasswordMismatch, b.CheckPassword("fred", "fredpass"))

	var users []string
	users, err = b.GetUsers()
	assert.NoError(t, err)
	assert.Equal(t, []string{"fred"}, users)
}

func TestApiAccount(t *testing.T) {
	var err error
	var b *Bank
	os.Remove("_account.db")
	b, err = NewBank("sqlite3", "_account.db")
	require.NoError(t, err)
	defer func() {
		b.Close()
		os.Remove("_account.db")
	}()
	require.NoError(t, b.NewBank())

	assert.NoError(t, b.NewAccount("fred"))
	assert.NoError(t, b.NewAccount("bob"))
	assert.NoError(t, b.NewAccount("house"))

	var accounts []string
	accounts, err = b.GetAccounts()
	assert.NoError(t, err)
	assert.Equal(t, []string{"bob", "fred", "house"}, accounts)

	var transactions []Transaction
	transactions, err = b.GetTransactions(10, 0, 0)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(transactions))

	// Transactions
	assert.NoError(t, b.NewTransaction("xfred", "house", "fred", "Buy a bus", 1000))
	assert.NoError(t, b.NewTransaction("xfred", "house", "fred", "Buy another bus", 1000))
	assert.NoError(t, b.NewTransaction("xbob", "house", "bob", "Buy a car", 3001))
	transactions, err = b.GetTransactions(10, 0, 0)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(transactions))

	assert.Equal(t, 1, transactions[2].ID)
	assert.Equal(t, "xfred", transactions[2].User)
	assert.Equal(t, "house", transactions[2].Origin)
	assert.Equal(t, "fred", transactions[2].Destination)
	assert.Equal(t, "Buy a bus", transactions[2].Description)
	assert.Equal(t, 1000, transactions[2].Amount)
	assert.Equal(t, 1000, transactions[2].DestinationBalanceAfter)
	assert.Equal(t, -1000, transactions[2].OriginBalanceAfter)

	assert.Equal(t, 2, transactions[1].ID)
	assert.Equal(t, "xfred", transactions[1].User)
	assert.Equal(t, "house", transactions[1].Origin)
	assert.Equal(t, "fred", transactions[1].Destination)
	assert.Equal(t, "Buy another bus", transactions[1].Description)
	assert.Equal(t, 1000, transactions[1].Amount)
	assert.Equal(t, -2000, transactions[1].OriginBalanceAfter)
	assert.Equal(t, 2000, transactions[1].DestinationBalanceAfter)

	assert.Equal(t, 3, transactions[0].ID)
	assert.Equal(t, "xbob", transactions[0].User)
	assert.Equal(t, "house", transactions[0].Origin)
	assert.Equal(t, "bob", transactions[0].Destination)
	assert.Equal(t, "Buy a car", transactions[0].Description)
	assert.Equal(t, 3001, transactions[0].Amount)
	assert.Equal(t, -5001, transactions[0].OriginBalanceAfter)
	assert.Equal(t, 3001, transactions[0].DestinationBalanceAfter)

	// Limit, offset
	transactions, err = b.GetTransactions(1, 0, 0)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(transactions))
	assert.Equal(t, 3, transactions[0].ID)

	transactions, err = b.GetTransactions(1, 2, 0)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(transactions))
	assert.Equal(t, 1, transactions[0].ID)

	// After
	transactions, err = b.GetTransactions(0, 0, 2)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(transactions))
	assert.Equal(t, 3, transactions[0].ID)

	// Distribution
	assert.NoError(t, b.Distribute("xfred", "house", []string{"bob", "fred"}, "distribution"))
	transactions, err = b.GetTransactions(2, 0, 0)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(transactions))

	assert.Equal(t, 4, transactions[1].ID)
	assert.Equal(t, "bob", transactions[1].Origin)
	assert.Equal(t, "house", transactions[1].Destination)
	assert.Equal(t, "distribution", transactions[1].Description)
	assert.Equal(t, 2500, transactions[1].Amount)
	assert.Equal(t, 501, transactions[1].OriginBalanceAfter)
	assert.Equal(t, -2501, transactions[1].DestinationBalanceAfter)

	assert.Equal(t, 5, transactions[0].ID)
	assert.Equal(t, "fred", transactions[0].Origin)
	assert.Equal(t, "house", transactions[0].Destination)
	assert.Equal(t, "distribution", transactions[0].Description)
	assert.Equal(t, 2500, transactions[0].Amount)
	assert.Equal(t, -500, transactions[0].OriginBalanceAfter)
	assert.Equal(t, -1, transactions[0].DestinationBalanceAfter)

}
