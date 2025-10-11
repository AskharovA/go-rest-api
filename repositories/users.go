package repositories

import (
	"AskharovA/go-rest-api/models"
	"AskharovA/go-rest-api/utils"
	"database/sql"
	"errors"
)

type UserRepository interface {
	Save(user *models.User) error
	ValidateCredentials(user *models.User) error
}

type dbUserRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) UserRepository {
	return &dbUserRepository{
		db: db,
	}
}

func (r *dbUserRepository) Save(user *models.User) error {
	query := `
	INSERT INTO users (email, password)
	VALUES (?, ?)
	`
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(user.Email, hashedPassword)
	if err != nil {
		return err
	}

	userId, err := result.LastInsertId()
	user.ID = userId
	return err
}

func (r *dbUserRepository) ValidateCredentials(user *models.User) error {
	query := "SELECT id, password FROM users WHERE email = ?"
	row := r.db.QueryRow(query, user.Email)

	var retrievedPassword string
	err := row.Scan(&user.ID, &retrievedPassword)
	if err != nil {
		return errors.New("credentials invalid")
	}

	passwordIsValid := utils.CheckPasswordHash(u.Password, retrievedPassword)
	if !passwordIsValid {
		return errors.New("credentials invalid")
	}

	return nil
}
