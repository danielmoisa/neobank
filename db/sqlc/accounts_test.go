package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/danielmoisa/neobank/utils"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)

	args := CreateAccountParams{
		Owner:    user.Username,
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, args.Owner, account.Owner)
	require.Equal(t, args.Balance, account.Balance)
	require.Equal(t, args.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestGetAccount(t *testing.T) {
	createAcc := createRandomAccount(t)
	getAcc, err := testQueries.GetAccount(context.Background(), createAcc.ID)

	require.NoError(t, err)
	require.NotEmpty(t, getAcc)
	require.Equal(t, createAcc.ID, getAcc.ID)
	require.Equal(t, createAcc.Owner, getAcc.Owner)
	require.Equal(t, createAcc.Balance, getAcc.Balance)
	require.Equal(t, createAcc.Currency, getAcc.Currency)
	require.WithinDuration(t, createAcc.CreatedAt, getAcc.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	acc := createRandomAccount(t)
	args := UpdateAccountParams{
		ID:      acc.ID,
		Balance: utils.RandomMoney(),
	}
	update, err := testQueries.UpdateAccount(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, update)
	require.Equal(t, acc.ID, update.ID)
	require.Equal(t, acc.Owner, update.Owner)
	require.Equal(t, args.Balance, update.Balance)
	require.Equal(t, acc.Currency, update.Currency)
	require.WithinDuration(t, acc.CreatedAt, update.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	acc := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), acc.ID)
	require.NoError(t, err)

	deletedAcc, err := testQueries.GetAccount(context.Background(), acc.ID)
	require.Error(t, err)
	require.Empty(t, deletedAcc)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
