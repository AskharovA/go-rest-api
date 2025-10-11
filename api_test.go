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

	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) (*sql.DB, int64) {
	testDbFile := "api_test.db"
	dbConn, err := db.InitDB(testDbFile)
	if err != nil {
		t.Fatalf("Could not initialize test database: %v", err)
	}

	err = db.CreateTables(dbConn)
	if err != nil {
		t.Fatalf("Could not create tables for test database: %v", err)
	}

	testEvent := models.Event{
		Name:        "Test Event",
		Description: "A test event.",
		Location:    "Test Location",
		DateTime:    time.Now(),
	}
	testEvent.Save(dbConn)

	testUser := models.User{
		Email:    "test@example.com",
		Password: "test",
	}
	testUser.Save(dbConn)

	t.Cleanup(
		func() {
			dbConn.Close()
			os.Remove(testDbFile)
		})

	return dbConn, testEvent.ID
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

func TestGetEvent(t *testing.T) {
	dbConn, eventId := setupTestDB(t)
	router := setupRouter(dbConn)

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
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateEvent(t *testing.T) {
	dbConn, eventId := setupTestDB(t)
	router := setupRouter(dbConn)

	url := fmt.Sprintf("/events/%d", eventId)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteEvent(t *testing.T) {
	dbConn, eventId := setupTestDB(t)
	router := setupRouter(dbConn)

	url := fmt.Sprintf("/events/%d", eventId)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
