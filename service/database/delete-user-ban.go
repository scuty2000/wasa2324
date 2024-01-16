package database

func (db *appdbimpl) DeleteUserBan(uuid string, bannedUUID string) error {
	_, err := db.c.Exec("DELETE FROM Bans WHERE UUID = ? AND BANNED_UUID = ?;", uuid, bannedUUID)
	return err
}
