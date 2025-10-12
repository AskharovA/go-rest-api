package repositories

import (
	"AskharovA/go-rest-api/models"
	"database/sql"
	"log"
)

type EventRepository interface {
	Save(event *models.Event) error
	GetAllEvents() ([]models.Event, error)
	GetEventByID(id int64) (*models.Event, error)
	Update(event *models.Event) error
	Delete(event *models.Event) error
	Register(userId int64, event *models.Event) error
	CancelRegistration(userId int64, event *models.Event) error
}

type dbEventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) EventRepository {
	return &dbEventRepository{
		db: db,
	}
}

func (r *dbEventRepository) Save(event *models.Event) error {
	query := `
	INSERT INTO events (name, description, location, dateTime, userId)
	VALUES (?, ?, ?, ?, ?)
	`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("Warning: could not close stmt: %v", err)
		}
	}()

	result, err := stmt.Exec(event.Name, event.Description, event.Location, event.DateTime, event.UserID)
	if err != nil {
		return err
	}

	eventId, err := result.LastInsertId()
	event.ID = eventId

	return err
}

func (r *dbEventRepository) GetAllEvents() ([]models.Event, error) {
	query := "SELECT * FROM events ORDER BY id"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Warning: could not close stmt: %v", err)
		}
	}()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)

		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (r *dbEventRepository) GetEventByID(id int64) (*models.Event, error) {
	query := "SELECT * FROM events WHERE id = ?"
	row := r.db.QueryRow(query, id)

	var event models.Event
	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)

	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *dbEventRepository) Update(event *models.Event) error {
	query := `
	UPDATE events
	SET name = ?, description = ?, location = ?, dateTime = ?
	WHERE id = ?
	`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("Warning: could not close stmt: %v", err)
		}
	}()

	_, err = stmt.Exec(event.Name, event.Description, event.Location, event.DateTime, event.ID)
	return err
}

func (r *dbEventRepository) Delete(event *models.Event) error {
	query := `
	DELETE FROM events
	WHERE id = ?
	`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("Warning: could not close stmt: %v", err)
		}
	}()

	_, err = stmt.Exec(event.ID)
	return err
}

func (r *dbEventRepository) Register(userId int64, event *models.Event) error {
	query := "INSERT INTO registrations (eventId, userId) VALUES (?, ?)"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("Warning: could not close stmt: %v", err)
		}
	}()

	_, err = stmt.Exec(event.ID, userId)
	return err
}

func (r *dbEventRepository) CancelRegistration(userId int64, event *models.Event) error {
	query := "DELETE FROM registrations WHERE eventId = ? AND userId = ?"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("Warning: could not close stmt: %v", err)
		}
	}()

	_, err = stmt.Exec(event.ID, userId)
	return err
}
