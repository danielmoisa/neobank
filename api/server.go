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
	router := echo.New()

	// Register validator
	router.Validator = &CustomValidator{validator: validator.New()}

	// Middleware
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())

	// Routes
	router.GET("/swagger/*", echoSwagger.WrapHandler)
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	router.POST("/payments", server.createPayment)
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	server.router = router
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
