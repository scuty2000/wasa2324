package database

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
)

func (db *appdbimpl) CreateUser(username string) (string, error) {
	identifier := uuid.New().String()
	var err error
	for db.checkCollision(identifier) {
		identifier = uuid.New().String()
	}
	_, err = db.c.Exec("INSERT INTO Users (UUID, USERNAME) VALUES (?, ?)", identifier, username)
	return identifier, err
}

func (db *appdbimpl) checkCollision(identifier string) bool {
	err := db.c.QueryRow("SELECT UUID FROM Users WHERE UUID=?", identifier).Scan()
	return !errors.Is(err, sql.ErrNoRows)
}
