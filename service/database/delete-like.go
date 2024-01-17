package database

func (db *appdbimpl) DeleteUserLike(userUUID string, photoUUID string) (int, error) {
	result, err := db.c.Exec("DELETE FROM Likes WHERE USER_UUID = ? AND PHOTO_UUID = ?;", userUUID, photoUUID)
	affectedRows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(affectedRows), nil
}
