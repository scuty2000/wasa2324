package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) SetSession(uuid string, bearer string) error {
	_, err := db.c.Exec("DELETE FROM Auth WHERE UUID=?", uuid)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	}
	_, err = db.c.Exec("INSERT INTO Auth (UUID, BEARER_TOKEN) VALUES (?, ?)", uuid, bearer)
	return err
}
