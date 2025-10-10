package main

import (
	"AskharovA/go-rest-api/db"
	"AskharovA/go-rest-api/routes"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	server := gin.Default()
	routes.RegisterRoutes(server)

	return server
}

func main() {
	db.InitDB()

	server := setupRouter()
	server.Run(":8080")
}
