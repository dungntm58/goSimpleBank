package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/dungntm58/goSimpleBank/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry() (CreateEntryParams, Entry, error) {
	_, acc, err := createRandomAccount()
	if err != nil {
		return CreateEntryParams{}, Entry{}, err
	}
	arg := CreateEntryParams{
		AccountID: acc.ID,
		Amount:    util.RandomAmount(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	return arg, entry, err
}

func TestCreateEntry(t *testing.T) {
	arg, entry, err := createRandomEntry()
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
}

func TestGetEntry(t *testing.T) {
	_, entry, err := createRandomEntry()
	require.NoError(t, err)

	entry2, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry.ID, entry2.ID)
	require.Equal(t, entry.AccountID, entry2.AccountID)
	require.Equal(t, entry.Amount, entry2.Amount)
	require.WithinDuration(t, entry.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestDeleteEntry(t *testing.T) {
	_, entry, _err := createRandomEntry()
	require.NoError(t, _err)

	err := testQueries.DeleteEntry(context.Background(), entry.ID)
	require.NoError(t, err)

	entry2, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, entry2)
}

func TestListEntries(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomEntry()
	}

	arg := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	entrys, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entrys, 5)

	for _, entry := range entrys {
		require.NotEmpty(t, entry)
	}
}
