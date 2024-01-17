package database

func (db *appdbimpl) SetUserLike(userUUID string, photoUUID string) error {
	_, err := db.c.Exec("INSERT INTO Likes (USER_UUID, PHOTO_UUID) VALUES (?, ?)", userUUID, photoUUID)
	return err
}
