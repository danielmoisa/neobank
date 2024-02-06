package api

import (
	"net/http"

	db "github.com/danielmoisa/neobank/db/sqlc"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	store  db.Store
	router *echo.Echo
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := echo.New()

	// Register validator
	router.Validator = &CustomValidator{validator: validator.New()}

	// Middleware
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())

	// Routes
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	server.router = router
	return server
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
