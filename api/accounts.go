package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	db "github.com/danielmoisa/neobank/db/sqlc"
	"github.com/danielmoisa/neobank/tokens"
	"github.com/labstack/echo/v4"
)

type createAccountRequest struct {
	Currency string `json:"currency" validate:"required"`
}

// createAccount godoc
// @Summary Create an account
// @Description Create a new account with the specified owner and currency.
// @Tags Accounts
// @Accept json
// @Produce json
// @Param request body createAccountRequest true "Request body for creating an account"
// @Success 201 {object} db.Account
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /accounts [post]
func (server *Server) createAccount(ctx echo.Context) error {
	req := new(createAccountRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	authPayload := ctx.Get(authorizationPayloadKey).(*tokens.Payload)

	account, err := server.store.CreateAccount(ctx.Request().Context(), db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
	}

	return ctx.JSON(http.StatusCreated, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" validate:"required"`
}

// getAccount godoc
// @Summary Get an account by ID
// @Description Retrieve an account by its unique ID.
// @Tags Accounts
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Success 200 {object} db.Account
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Account Not Found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /accounts/{id} [get]
func (server *Server) getAccount(ctx echo.Context) error {
	req := new(getAccountRequest)

	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID"})
	}

	req.ID = id

	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	account, err := server.store.GetAccount(ctx.Request().Context(), req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, ErrorResponse{Error: "Account not found"})
		}

		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal server error"})
	}

	authPayload := ctx.Get(authorizationPayloadKey).(*tokens.Payload)
	if account.Owner != authPayload.Username {
		err := errors.New("account doesn't belongs to auth user")
		return ctx.JSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
	}

	return ctx.JSON(http.StatusOK, account)
}

type listAccountRequest struct {
	PageID   int32 `form:"page_id" validate:"required,min=1"`
	PageSize int32 `form:"page_size" validate:"required,min=5,max=10"`
}

// listAccounts godoc
// @Summary List accounts
// @Description Get a list of accounts with pagination.
// @Tags Accounts
// @Accept json
// @Produce json
// @Param page_id query int true "Page ID for pagination"
// @Param page_size query int true "Number of accounts per page (min: 5, max: 10)"
// @Success 200 {array} db.Account
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /accounts [get]
func (server *Server) listAccounts(ctx echo.Context) error {
	req := new(listAccountRequest)

	pageIDStr := ctx.QueryParam("page_id")
	pageSizeStr := ctx.QueryParam("page_size")

	pageID, err := strconv.ParseInt(pageIDStr, 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid page_id"})
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid page_size"})
	}

	req.PageID = int32(pageID)
	req.PageSize = int32(pageSize)

	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	authPayload := ctx.Get(authorizationPayloadKey).(*tokens.Payload)

	accounts, err := server.store.ListAccounts(ctx.Request().Context(), db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal server error"})
	}

	return ctx.JSON(http.StatusOK, accounts)
}
