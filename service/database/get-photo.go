package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) GetPhoto(photoUUID string, requestingUUID string) (string, string, int, int, bool, error) {
	var ownerUUID string
	var date string
	err := db.c.QueryRow("SELECT OWNER_UUID, DATE FROM Photos WHERE UUID=?", photoUUID).Scan(&ownerUUID, &date)
	if err != nil {
		return "", "", 0, 0, false, err
	}
	var likes int
	err = db.c.QueryRow("SELECT COUNT(*) FROM Likes WHERE PHOTO_UUID=?", photoUUID).Scan(&likes)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			likes = 0
		} else {
			return "", "", 0, 0, false, err
		}
	}
	var comments int
	err = db.c.QueryRow("SELECT COUNT(*) FROM Comments WHERE PHOTO_UUID=?", photoUUID).Scan(&comments)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			comments = 0
		} else {
			return "", "", 0, 0, false, err
		}
	}
	var isLiked int
	err = db.c.QueryRow("SELECT COUNT(*) FROM Likes WHERE PHOTO_UUID=? AND USER_UUID=?", photoUUID, requestingUUID).Scan(&isLiked)
	return ownerUUID, date, likes, comments, isLiked == 1, err
}
