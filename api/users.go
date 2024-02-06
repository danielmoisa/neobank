package api

import (
	"net/http"
	"time"

	db "github.com/danielmoisa/neobank/db/sqlc"
	"github.com/danielmoisa/neobank/utils"
	"github.com/labstack/echo/v4"
)

type createUserRequest struct {
	Username string `json:"username" validate:"required,alphanum"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}
type createUserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

// createUser godoc
// @Summary Create a user
// @Description Create a new user with the specified details.
// @Tags Users
// @Accept json
// @Produce json
// @Param request body createUserRequest true "Request body for creating a user"
// @Success 201 {object} createUserResponse
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /users [post]
func (server *Server) createUser(ctx echo.Context) error {
	req := new(createUserRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
	}

	user, err := server.store.CreateUser(ctx.Request().Context(), db.CreateUserParams{
		Username:       req.Username,
		FullName:       req.FullName,
		Email:          req.Email,
		HashedPassword: hash,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	res := createUserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
	return ctx.JSON(http.StatusCreated, res)
}
