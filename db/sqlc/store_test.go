package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPaymentTx(t *testing.T) {
	store := NewStore(testDB)

	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	n := 5
	amount := int64(10)

	// run n concurrent transfer transaction
	errs := make(chan error)
	results := make(chan PaymentTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.PaymentTx(context.Background(), PaymentTxParams{
				FromAccountID: acc1.ID,
				ToAccountID:   acc2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	// check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		payment := result.Payment
		require.NotEmpty(t, payment)
		require.Equal(t, acc1.ID, payment.FromAccountID)
		require.Equal(t, acc2.ID, payment.ToAccountID)
		require.Equal(t, amount, payment.Amount)
		require.NotZero(t, payment.ID)
		require.NotZero(t, payment.CreatedAt)

		_, err = store.GetPayment(context.Background(), payment.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, acc1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, acc2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts' balance
	}

}