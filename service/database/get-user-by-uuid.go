package database

func (db *appdbimpl) GetUserByUUID(uuid string) (string, error) {
	var name string
	err := db.c.QueryRow("SELECT USERNAME FROM Users WHERE UUID=?", uuid).Scan(&name)
	return name, err
}
