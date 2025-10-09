package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

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
	server := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/events", nil)

	server.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
