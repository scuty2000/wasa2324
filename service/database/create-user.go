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
	err = db.c.QueryRow("INSERT INTO Users (UUID, USERNAME) VALUES (?, ?)", identifier, username).Scan()
	return identifier, err
}

func (db *appdbimpl) checkCollision(identifier string) bool {
	var err error
	err = db.c.QueryRow("SELECT UUID FROM Users WHERE UUID=?", identifier).Scan()
	return !errors.Is(err, sql.ErrNoRows)
}
