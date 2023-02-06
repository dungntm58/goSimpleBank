package db

import (
	"context"
	"testing"

	"github.com/dungntm58/goSimpleBank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccountWithFixedBalance(balance int64, currency string) (CreateAccountParams, Account, error) {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  balance,
		Currency: currency,
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	return arg, account, err
}

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	balance := int64(300)
	currency := util.RandomCurrency()

	_, acc1, err := createRandomAccountWithFixedBalance(balance, currency)
	require.NoError(t, err)

	_, acc2, err := createRandomAccountWithFixedBalance(balance, currency)
	require.NoError(t, err)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			res, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: acc1.ID,
				ToAccountID:   acc2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- res
		}()
	}

	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		res := <-results
		require.NotEmpty(t, res)

		tx := res.Tx
		require.NotEmpty(t, tx)
		require.Equal(t, acc1.ID, tx.FromAccountID)
		require.Equal(t, acc2.ID, tx.ToAccountID)
		require.Equal(t, amount, tx.Amount)
		require.NotZero(t, tx.ID)
		require.NotZero(t, tx.CreatedAt)

		_, err = store.GetTransaction(context.Background(), tx.ID)
		require.NoError(t, err)

		fromEntry := res.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, acc1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := res.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, acc2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := res.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, acc1.ID)

		toAccount := res.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, acc2.ID)

		diff1 := acc1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - acc2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final accounts' balances
	updatedAcc1, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	updatedAcc2, err := testQueries.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)

	require.Equal(t, acc1.Balance-int64(n)*amount, updatedAcc1.Balance)
	require.Equal(t, acc2.Balance+int64(n)*amount, updatedAcc2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	balance := int64(300)
	currency := util.RandomCurrency()

	_, acc1, err := createRandomAccountWithFixedBalance(balance, currency)
	require.NoError(t, err)

	_, acc2, err := createRandomAccountWithFixedBalance(balance, currency)
	require.NoError(t, err)

	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		isEven := i%2 == 0
		go func() {
			var fromAccountID, toAccountID int64
			if isEven {
				fromAccountID = acc1.ID
				toAccountID = acc2.ID
			} else {
				fromAccountID = acc2.ID
				toAccountID = acc1.ID
			}
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err, i)
	}

	// check the final accounts' balances
	updatedAcc1, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	updatedAcc2, err := testQueries.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)

	require.Equal(t, acc1.Balance, updatedAcc1.Balance)
	require.Equal(t, acc2.Balance, updatedAcc2.Balance)
}
