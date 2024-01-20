package database

func (db *appdbimpl) DeleteUserFollow(uuid string, followedUUID string) error {
	_, err := db.c.Exec("DELETE FROM Follows WHERE UUID = ? AND FOLLOWED_UUID = ?;", uuid, followedUUID)
	return err
}
