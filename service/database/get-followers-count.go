package database

func (db *appdbimpl) GetFollowersCount(uuid string) (int, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(1) FROM Follows WHERE FOLLOWED_UUID=?", uuid).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
