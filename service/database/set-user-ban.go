package database

func (db *appdbimpl) SetUserBan(uuid string, bannedUUID string) error {
	_, err := db.c.Exec("INSERT INTO Bans (UUID, BANNED_UUID) VALUES (?, ?)", uuid, bannedUUID)
	return err
}
