package models

import (
	"AskharovA/go-rest-api/db"
	"time"
)

type Event struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Location    string    `json:"location" binding:"required"`
	DateTime    time.Time `json:"dateTime" binding:"required"`
	UserID      int       `json:"userId"`
}

var events = []Event{}

func (e *Event) Save() error {
	query := `
	INSERT INTO events (name, description, location, dateTime, userId)
	VALUES (?, ?, ?, ?, ?)
	`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()
	result, err := stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserID)

	if err != nil {
		return err
	}

	eventId, err := result.LastInsertId()
	e.ID = eventId

	return err
}

func GetAllEvents() []Event {
	return events
}
