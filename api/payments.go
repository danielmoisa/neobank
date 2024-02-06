package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/danielmoisa/neobank/db/sqlc"
	"github.com/labstack/echo/v4"
)

type paymentRequest struct {
	FromAccountID int64  `json:"from_account_id" validation:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" validation:"required,min=1"`
	Amount        int64  `json:"amount" validation:"required,gt=0"`
	Currency      string `json:"currency" validation:"required,oneof=USD EUR CAD"`
}

func (server *Server) createPayment(ctx echo.Context) error {
	req := new(paymentRequest)

	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if !server.validAccount(ctx, req.FromAccountID, req.Currency) {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid account"})
	}

	if !server.validAccount(ctx, req.ToAccountID, req.Currency) {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid account"})
	}

	payment, err := server.store.PaymentTx(ctx.Request().Context(), db.PaymentTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	})

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})

	}

	return ctx.JSON(http.StatusCreated, payment)

}

func (server *Server) validAccount(ctx echo.Context, accountID int64, currency string) bool {
	account, err := server.store.GetAccount(ctx.Request().Context(), accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
			return false
		}

		ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return false
	}

	return true
}
