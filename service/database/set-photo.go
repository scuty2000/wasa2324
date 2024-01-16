package database

import (
	"github.com/google/uuid"
	"time"
)

func (db *appdbimpl) SetPhoto(ownerUUID string) (string, error) {
	photoUUID := uuid.New().String()
	currentTime := time.Now().Format("2006-01-02 15:04:05.000000")
	_, err := db.c.Exec("INSERT INTO Photos (UUID, OWNER_UUID, DATE) VALUES (?, ?, ?)", photoUUID, ownerUUID, currentTime)
	if err != nil {
		for err.Error() == "UNIQUE constraint failed: Photos.UUID" {
			photoUUID = uuid.New().String()
			_, err = db.c.Exec("INSERT INTO Photos (UUID, OWNER_UUID, DATE) VALUES (?, ?, ?)", photoUUID, ownerUUID, currentTime)
		}
		if err != nil {
			return "", err
		}
	}
	return photoUUID, err
}
