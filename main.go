package main

import (
	"AskharovA/go-rest-api/config"
	"AskharovA/go-rest-api/db"
	"AskharovA/go-rest-api/routes"
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
)

func setupRouter(dbConn *sql.DB) *gin.Engine {
	server := gin.Default()
	routes.RegisterRoutes(server, dbConn)

	return server
}

func main() {
	dbConn, err := db.InitDB("api.db")
	if err != nil {
		panic("Could not connect to database.")
	}
	defer dbConn.Close()

	err = db.CreateTables(dbConn)
	if err != nil {
		panic("Could not create tables.")
	}

	server := setupRouter(dbConn)

	addr := fmt.Sprintf(":%d", config.Get().AppPort)
	server.Run(addr)
}
