package database

func (db *appdbimpl) GetUserByName(username string) (string, error) {
	var name string
	err := db.c.QueryRow("SELECT UUID FROM Users WHERE USERNAME=?", username).Scan(&name)
	return name, err
}
