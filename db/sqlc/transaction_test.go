package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/dungntm58/goSimpleBank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransaction() (CreateTransactionParams, Transaction, error) {
	_, acc1, err := createRandomAccount()
	if err != nil {
		return CreateTransactionParams{}, Transaction{}, err
	}
	_, acc2, err := createRandomAccount()
	if err != nil {
		return CreateTransactionParams{}, Transaction{}, err
	}

	arg := CreateTransactionParams{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Amount:        util.RandomAmount(),
	}

	tx, err := testQueries.CreateTransaction(context.Background(), arg)
	return arg, tx, err
}

func TestCreateTransaction(t *testing.T) {
	arg, tx, err := createRandomTransaction()
	require.NoError(t, err)
	require.NotEmpty(t, tx)

	require.Equal(t, arg.FromAccountID, tx.FromAccountID)
	require.Equal(t, arg.ToAccountID, tx.ToAccountID)
	require.Equal(t, arg.Amount, tx.Amount)

	require.NotZero(t, tx.ID)
	require.NotZero(t, tx.CreatedAt)
}

func TestGetTransaction(t *testing.T) {
	_, tx, err := createRandomTransaction()
	require.NoError(t, err)

	tx2, err := testQueries.GetTransaction(context.Background(), tx.ID)
	require.NoError(t, err)
	require.NotEmpty(t, tx2)

	require.Equal(t, tx.ID, tx2.ID)
	require.Equal(t, tx.ToAccountID, tx2.ToAccountID)
	require.Equal(t, tx.Amount, tx2.Amount)
	require.WithinDuration(t, tx.CreatedAt, tx2.CreatedAt, time.Second)
}

func TestDeleteTransaction(t *testing.T) {
	_, tx, _err := createRandomTransaction()
	require.NoError(t, _err)

	err := testQueries.DeleteTransaction(context.Background(), tx.ID)
	require.NoError(t, err)

	tx2, err := testQueries.GetTransaction(context.Background(), tx.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, tx2)
}

func TestListTransactions(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomTransaction()
	}

	arg := ListTransactionsParams{
		Limit:  5,
		Offset: 5,
	}

	txs, err := testQueries.ListTransactions(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, txs, 5)

	for _, tx := range txs {
		require.NotEmpty(t, tx)
	}
}
