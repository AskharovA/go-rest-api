package routes

import (
	"AskharovA/go-rest-api/models"
	"AskharovA/go-rest-api/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EventAPI struct {
	EventService *services.EventService
}

func (api *EventAPI) getEvents(context *gin.Context) {
	page := context.Query("page")
	per_page := context.Query("per_page")

	pageInt := 1
	perPageInt := 10

	if page != "" && per_page != "" {
		pageInt, _ = strconv.Atoi(page)
		perPageInt, _ = strconv.Atoi(per_page)
	}

	events, err := api.EventService.GetEvents(pageInt, perPageInt)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch events. Try again later."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "OK", "data": events})
}

func (api *EventAPI) getEvent(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event id."})
		return
	}

	event, err := api.EventService.GetEvent(eventId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not fetch event."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "OK", "data": event})
}

func (api *EventAPI) createEvent(context *gin.Context) {
	var event models.Event
	err := context.ShouldBindJSON(&event)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse requrest data."})
		return
	}

	userId := context.GetInt64("userId")
	event.UserID = userId
	err = api.EventService.CreateEvent(&event)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create event. Try again later."})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "OK", "data": event})
}

func (api *EventAPI) updateEvent(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event id."})
		return
	}

	userId := context.GetInt64("userId")
	event, err := api.EventService.GetEvent(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch the event."})
		return
	}

	if event.UserID != userId {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Not authorized to update event."})
		return
	}

	var updatedEvent models.Event
	err = context.ShouldBindJSON(&updatedEvent)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse requrest data."})
		return
	}

	updatedEvent.ID = eventId
	err = api.EventService.UpdateEvent(&updatedEvent)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update event."})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Event updated successfully!"})
}

func (api *EventAPI) deleteEvent(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse event id."})
		return
	}

	userId := context.GetInt64("userId")
	event, err := api.EventService.GetEvent(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch the event."})
		return
	}

	if event.UserID != userId {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Not authorized to delete event."})
		return
	}

	err = api.EventService.DeleteEvent(event)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete the event."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully!"})
}
