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

// createUser godoc
// @Summary Create a user
// @Description Create a new user with the specified details.
// @Tags Users
// @Accept json
// @Produce json
// @Param request body createUserRequest true "Request body for creating a user"
// @Success 201 {object} userResponse
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

	res := newUserResponse(user)
	return ctx.JSON(http.StatusCreated, res)
}

type loginUserRequest struct {
	Username string `json:"username" validate:"required,alphanum"`
	Password string `json:"password" validate:"required,min=6"`
}

type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

// createUser godoc
// @Summary Login a user
// @Description Login a new user with the specified details.
// @Tags Users
// @Accept json
// @Produce json
// @Param request body loginUserRequest true "Request body for login in a user"
// @Success 201 {object} loginUserResponse
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /users/login [post]
func (server *Server) loginUser(ctx echo.Context) error {
	req := new(loginUserRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	if err := ctx.Validate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	user, err := server.store.GetUser(ctx.Request().Context(), req.Username)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Invalid credentials"})
	}

	err = utils.CheckPassword(user.HashedPassword, req.Password)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid credentials"})
	}

	token, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
	}

	res := loginUserResponse{
		AccessToken: token,
		User:        newUserResponse(user),
	}

	return ctx.JSON(http.StatusOK, res)
}
