package main

import (
	"AskharovA/go-rest-api/db"
	"AskharovA/go-rest-api/models"
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *sql.DB {
	testDbFile := "api_test.db"
	dbConn, err := db.InitDB(testDbFile)
	if err != nil {
		t.Fatalf("Could not initialize test database: %v", err)
	}

	err = db.CreateTables(dbConn)
	if err != nil {
		t.Fatalf("Could not create tables for test database: %v", err)
	}

	t.Cleanup(
		func() {
			dbConn.Close()
			os.Remove(testDbFile)
		})

	return dbConn
}

func TestGetEvents(t *testing.T) {
	dbConn := setupTestDB(t)
	router := setupRouter(dbConn)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/events", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateEvent(t *testing.T) {
	dbConn := setupTestDB(t)
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
