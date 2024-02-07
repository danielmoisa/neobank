package main

import (
	"database/sql"
	"log"

	"github.com/danielmoisa/neobank/api"
	db "github.com/danielmoisa/neobank/db/sqlc"
	"github.com/danielmoisa/neobank/utils"
	_ "github.com/lib/pq"
)

// @title Swagger Neobank API
// @version 1.0
// @description This is a sample bank server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host neobank.swagger.io
func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}
