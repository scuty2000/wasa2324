package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) DeletePhoto(photoUUID string) error {
	_, err := db.c.Exec("DELETE FROM Photos WHERE UUID = ?;", photoUUID)
	if err != nil {
		return err
	}
	_, err = db.c.Exec("DELETE FROM Likes WHERE PHOTO_UUID = ?;", photoUUID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	}
	_, err = db.c.Exec("DELETE FROM Comments WHERE PHOTO_UUID = ?;", photoUUID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	}
	return err
}
