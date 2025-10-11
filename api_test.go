package main

import (
	"AskharovA/go-rest-api/db"
	"AskharovA/go-rest-api/models"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) (*sql.DB, int64) {
	testDbFile := "api_test.db"
	os.Remove(testDbFile)

	dbConn, err := db.InitDB(testDbFile)
	if err != nil {
		t.Fatalf("Could not initialize test database: %v", err)
	}

	err = db.CreateTables(dbConn)
	if err != nil {
		t.Fatalf("Could not create tables for test database: %v", err)
	}

	testUser := models.User{
		Email:    "test@example.com",
		Password: "test",
	}
	testUser.Save(dbConn)

	testEvent := models.Event{
		Name:        "Test Event",
		Description: "A test event.",
		Location:    "Test Location",
		DateTime:    time.Now(),
		UserID:      testUser.ID,
	}
	testEvent.Save(dbConn)

	t.Cleanup(
		func() {
			dbConn.Close()
			os.Remove(testDbFile)
		})

	return dbConn, testEvent.ID
}

func getAuthToken(t *testing.T, router *gin.Engine) string {
	loginCredentials := gin.H{
		"email":    "test@example.com",
		"password": "test",
	}
	payload, _ := json.Marshal(loginCredentials)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(string(payload)))

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var responseBody map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.Nil(t, err)

	token, exists := responseBody["token"]
	assert.True(t, exists)
	assert.NotEmpty(t, token)

	return token
}

func TestGetEvents(t *testing.T) {
	dbConn, _ := setupTestDB(t)
	router := setupRouter(dbConn)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/events", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateEvent(t *testing.T) {
	dbConn, _ := setupTestDB(t)
	router := setupRouter(dbConn)
	token := getAuthToken(t, router)

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
	req.Header.Set("Authorization", token)

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

func TestGetEvent(t *testing.T) {
	dbConn, eventId := setupTestDB(t)
	router := setupRouter(dbConn)

	url := fmt.Sprintf("/events/%d", eventId)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateEvent(t *testing.T) {
	dbConn, eventId := setupTestDB(t)
	router := setupRouter(dbConn)
	token := getAuthToken(t, router)

	updatedEvent := models.Event{
		Name:        "Test Event",
		Description: "A test event.",
		Location:    "Test Location",
		DateTime:    time.Now(),
	}
	payload, _ := json.Marshal(updatedEvent)

	url := fmt.Sprintf("/events/%d", eventId)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteEvent(t *testing.T) {
	dbConn, eventId := setupTestDB(t)
	router := setupRouter(dbConn)
	token := getAuthToken(t, router)

	url := fmt.Sprintf("/events/%d", eventId)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
