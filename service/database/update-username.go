package database

func (db *appdbimpl) UpdateUsername(uuid string, username string) error {
	_, err := db.c.Exec("UPDATE Users SET USERNAME = ? WHERE UUID = ?", username, uuid)
	return err
}
