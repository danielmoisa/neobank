package db

import (
	"context"
	"testing"
	"time"

	"github.com/danielmoisa/neobank/utils"
	"github.com/stretchr/testify/require"
)

func createRandomPayment(t *testing.T, account1, account2 Account) Payment {
	args := CreatePaymentParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        utils.RandomMoney(),
	}

	payment, err := testQueries.CreatePayment(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, payment)

	require.Equal(t, args.FromAccountID, payment.FromAccountID)
	require.Equal(t, args.ToAccountID, payment.ToAccountID)
	require.Equal(t, args.Amount, payment.Amount)

	require.NotZero(t, payment.ID)
	require.NotZero(t, payment.CreatedAt)

	return payment
}

func TestCreatePayment(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	createRandomPayment(t, account1, account2)
}

func TestGetPayment(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	payment1 := createRandomPayment(t, account1, account2)

	payment2, err := testQueries.GetPayment(context.Background(), payment1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, payment2)

	require.Equal(t, payment1.ID, payment2.ID)
	require.Equal(t, payment1.FromAccountID, payment2.FromAccountID)
	require.Equal(t, payment1.ToAccountID, payment2.ToAccountID)
	require.Equal(t, payment1.Amount, payment2.Amount)
	require.WithinDuration(t, payment1.CreatedAt, payment2.CreatedAt, time.Second)
}

func TestListPayment(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	for i := 0; i < 5; i++ {
		createRandomPayment(t, account1, account2)
		createRandomPayment(t, account2, account1)
	}

	arg := ListPaymentsParams{
		FromAccountID: account1.ID,
		ToAccountID:   account1.ID,
		Limit:         5,
		Offset:        5,
	}

	payments, err := testQueries.ListPayments(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, payments, 5)

	for _, payment := range payments {
		require.NotEmpty(t, payment)
		require.True(t, payment.FromAccountID == account1.ID || payment.ToAccountID == account1.ID)
	}
}
