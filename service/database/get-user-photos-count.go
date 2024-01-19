package database

func (db *appdbimpl) GetUserPhotosCount(uuid string) (int, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(1) FROM Photos WHERE OWNER_UUID=?", uuid).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
