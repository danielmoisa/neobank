package api

import (
	"fmt"
	"net/http"

	db "github.com/danielmoisa/neobank/db/sqlc"
	_ "github.com/danielmoisa/neobank/docs"
	"github.com/danielmoisa/neobank/tokens"
	"github.com/danielmoisa/neobank/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Server struct {
	store      db.Store
	router     *echo.Echo
	tokenMaker tokens.Maker
	config     utils.Config
}

func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := tokens.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{store: store, tokenMaker: tokenMaker, config: config}
	e := echo.New()

	// Register validator
	e.Validator = &CustomValidator{validator: validator.New()}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/users", server.createUser)
	e.POST("/users/login", server.loginUser)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Protected routes
	e.POST("/accounts", server.createAccount, authMiddleware(server.tokenMaker))
	e.GET("/accounts/:id", server.getAccount, authMiddleware(server.tokenMaker))
	e.GET("/accounts", server.listAccounts, authMiddleware(server.tokenMaker))
	e.POST("/payments", server.createPayment, authMiddleware(server.tokenMaker))

	server.router = e
	return server, nil
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Start(address)
}

// CustomValidator represents a custom validator for Echo
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

type ErrorResponse struct {
	Error string `json:"error"`
}
