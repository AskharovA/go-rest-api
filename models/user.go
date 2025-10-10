package models

import (
	"AskharovA/go-rest-api/utils"
	"database/sql"
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
