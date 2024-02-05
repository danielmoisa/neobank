package api

import (
	"database/sql"
	"net/http"
	"strconv"

	db "github.com/danielmoisa/neobank/db/sqlc"
	"github.com/labstack/echo/v4"
)

type createAccountRequest struct {
	Owner    string `json:"owner" validated:"required"`
	Currency string `json:"currency" validate:"required"`
}

func (server *Server) createAccount(ctx echo.Context) error {
	req := new(createAccountRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	account, err := server.store.CreateAccount(ctx.Request().Context(), db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusCreated, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" validate:"required"`
}

func (server *Server) getAccount(ctx echo.Context) error {
	req := new(getAccountRequest)

	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	req.ID = id

	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	account, err := server.store.GetAccount(ctx.Request().Context(), req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Account not found"})
		}

		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return ctx.JSON(http.StatusOK, account)
}

type listAccountRequest struct {
	PageID   int32 `form:"page_id" validate:"required,min=1"`
	PageSize int32 `form:"page_size" validate:"required,min=5,max=10"`
}

func (server *Server) listAccounts(ctx echo.Context) error {
	req := new(listAccountRequest)

	pageIDStr := ctx.QueryParam("page_id")
	pageSizeStr := ctx.QueryParam("page_size")

	pageID, err := strconv.ParseInt(pageIDStr, 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid page_id"})
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid page_size"})
	}

	req.PageID = int32(pageID)
	req.PageSize = int32(pageSize)

	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	accounts, err := server.store.ListAccounts(ctx.Request().Context(), db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return ctx.JSON(http.StatusOK, accounts)
}
