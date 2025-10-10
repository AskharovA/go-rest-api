package main

import (
	"AskharovA/go-rest-api/db"
	"AskharovA/go-rest-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	server := gin.Default()

	server.GET("/events", getEvents)
	server.POST("/events", createEvent)

	server.Run(":8080")
}

func getEvents(context *gin.Context) {
	events := models.GetAllEvents()
	context.JSON(http.StatusOK, gin.H{"message": "OK", "data": events})
}

func createEvent(context *gin.Context) {
	var event models.Event
	err := context.ShouldBindJSON(&event)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse requrest data."})
		return
	}

	event.ID = 1     // dummy value
	event.UserID = 1 // dummy value
	event.Save()

	context.JSON(http.StatusCreated, gin.H{"message": "OK", "data": event})
}
