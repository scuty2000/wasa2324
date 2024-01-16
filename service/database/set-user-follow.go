package database

func (db *appdbimpl) SetUserFollow(uuid string, followedUUID string) error {
	_, err := db.c.Exec("INSERT INTO Follows (UUID, FOLLOWED_UUID) VALUES (?, ?)", uuid, followedUUID)
	return err
}
