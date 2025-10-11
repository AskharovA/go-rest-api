package routes

import (
	"AskharovA/go-rest-api/middlewares"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine, dbConn *sql.DB) {
	server.GET("/events", func(c *gin.Context) { getEvents(c, dbConn) })
	server.GET("/events/:id", func(c *gin.Context) { getEvent(c, dbConn) })

	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate)
	authenticated.POST("/events", func(c *gin.Context) { createEvent(c, dbConn) })
	authenticated.PUT("/events/:id", func(c *gin.Context) { updateEvent(c, dbConn) })
	authenticated.DELETE("/events/:id", func(c *gin.Context) { deleteEvent(c, dbConn) })
	authenticated.POST("/events/:id/register", func(c *gin.Context) { registerForEvent(c, dbConn) })
	authenticated.DELETE("/events/:id/register", func(c *gin.Context) { cancelRegistration(c, dbConn) })

	server.POST("/signup", func(c *gin.Context) { signup(c, dbConn) })
	server.POST("/login", func(c *gin.Context) { login(c, dbConn) })
}
