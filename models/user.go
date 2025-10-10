package models

import (
	"AskharovA/go-rest-api/utils"
	"database/sql"
	"errors"
)

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u *User) Save(dbConn *sql.DB) error {
	query := `
	INSERT INTO users (email, password)
	VALUES (?, ?)
	`
	stmt, err := dbConn.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(u.Email, hashedPassword)
	if err != nil {
		return err
	}

	userId, err := result.LastInsertId()
	u.ID = userId
	return err
}

func (u *User) ValidateCredentials(dbConn *sql.DB) error {
	query := "SELECT id, password FROM users WHERE email = ?"
	row := dbConn.QueryRow(query, u.Email)

	var retrievedPassword string
	err := row.Scan(&u.ID, &retrievedPassword)
	if err != nil {
		return errors.New("credentials invalid")
	}

	passwordIsValid := utils.CheckPasswordHash(u.Password, retrievedPassword)
	if !passwordIsValid {
		return errors.New("credentials invalid")
	}

	return nil
}
