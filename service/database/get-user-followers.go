package database

func (db *appdbimpl) GetUserFollowers(uuid string) ([]string, error) {
	var followers []string
	rows, err := db.c.Query("SELECT UUID FROM Follows WHERE FOLLOWED_UUID=?", uuid)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	defer rows.Close()
	for rows.Next() {
		var follower string
		err = rows.Scan(&follower)
		if err != nil {
			return nil, err
		}
		followers = append(followers, follower)
	}
	if err != nil {
		return nil, err
	}
	return followers, nil
}
