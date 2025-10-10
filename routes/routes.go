package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine, dbConn *sql.DB) {
	server.GET("/events", func(c *gin.Context) { getEvents(c, dbConn) })
	server.GET("/events/:id", func(c *gin.Context) { getEvent(c, dbConn) })
	server.POST("/events", func(c *gin.Context) { createEvent(c, dbConn) })
	server.PUT("/events/:id", func(c *gin.Context) { updateEvent(c, dbConn) })
}
