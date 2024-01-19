package database

func (db *appdbimpl) GetFollowingCount(uuid string) (int, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(1) FROM Follows WHERE UUID=?", uuid).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
