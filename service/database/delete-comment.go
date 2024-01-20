package database

func (db *appdbimpl) DeleteComment(photoUUID string, commentUUID string) (int, error) {
	result, err := db.c.Exec("DELETE FROM Comments WHERE PHOTO_UUID = ? AND COMMENT_UUID;", photoUUID, commentUUID)
	if err != nil {
		return 0, err
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(affectedRows), nil
}
