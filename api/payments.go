package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/danielmoisa/neobank/db/sqlc"
	"github.com/danielmoisa/neobank/tokens"
	"github.com/labstack/echo/v4"
)

type paymentRequest struct {
	FromAccountID int64  `json:"from_account_id" validation:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" validation:"required,min=1"`
	Amount        int64  `json:"amount" validation:"required,gt=0"`
	Currency      string `json:"currency" validation:"required,oneof=USD EUR CAD"`
}

// createPayment godoc
// @Summary Create a payment
// @Description Transfer funds between two accounts.
// @Tags Payments
// @Accept json
// @Produce json
// @Param request body paymentRequest true "Request body for creating a payment"
// @Success 201 {object} db.Payment
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Account Not Found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /payments [post]
func (server *Server) createPayment(ctx echo.Context) error {
	req := new(paymentRequest)

	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	fromAccount, valid := server.validAccount(ctx, req.FromAccountID, req.Currency)

	if !valid {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid account"})
	}

	authPayload := ctx.Get(authorizationPayloadKey).(*tokens.Payload)
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("account doesn't belongs to auth user")
		return ctx.JSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
	}

	_, valid = server.validAccount(ctx, req.ToAccountID, req.Currency)

	if !valid {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid account"})
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

func (server *Server) validAccount(ctx echo.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx.Request().Context(), accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
			return account, false
		}

		ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return account, false
	}

	return account, true
}
