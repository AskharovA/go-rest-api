package main

import (
	"AskharovA/go-rest-api/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/events", getEvents)
	r.POST("/events", createEvent)
	return r
}

func TestGetEvents(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/events", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateEvent(t *testing.T) {
	router := setupRouter()

	newEvent := models.Event{
		Name:        "Test Event",
		Description: "A test event.",
		Location:    "Test Location",
		DateTime:    time.Now(),
	}
	payload, _ := json.Marshal(newEvent)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/events", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var responseBody map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.Nil(t, err)

	eventData, ok := responseBody["data"].(map[string]interface{})
	assert.True(t, ok)

	assert.Equal(t, newEvent.Name, eventData["name"])
	assert.Equal(t, newEvent.Description, eventData["description"])
}
