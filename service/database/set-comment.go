package database

import (
	"github.com/google/uuid"
	"time"
)

func (db *appdbimpl) SetComment(ownerUUID string, photoUUID string, commentText string) (string, string, error) {
	commentUUID := uuid.New().String()
	currentTime := time.Now().Format("2006-01-02 15:04:05.000000")
	_, err := db.c.Exec("INSERT INTO Comments (COMMENT_UUID, OWNER_UUID, PHOTO_UUID, DATE, COMMENT_TEXT) VALUES (?, ?, ?, ?, ?)", commentUUID, ownerUUID, photoUUID, currentTime, commentText)
	if err != nil {
		for err.Error() == "UNIQUE constraint failed: Comments.COMMENT_UUID" {
			commentUUID = uuid.New().String()
			_, err = db.c.Exec("INSERT INTO Comments (COMMENT_UUID, OWNER_UUID, PHOTO_UUID, DATE, COMMENT_TEXT) VALUES (?, ?, ?, ?, ?)", commentUUID, ownerUUID, photoUUID, currentTime, commentText)
		}
		if err != nil {
			return "", "", err
		}
	}
	return commentUUID, currentTime, err
}
