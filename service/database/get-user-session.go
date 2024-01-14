package database

func (db *appdbimpl) GetUserSession(uuid string) (string, error) {
	var bearer string
	err := db.c.QueryRow("SELECT BEARER_TOKEN FROM Auth WHERE UUID = ?", uuid).Scan(&bearer)
	return bearer, err
}
