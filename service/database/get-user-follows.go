package database

func (db *appdbimpl) GetUserFollows(uuid string) ([]string, error) {
	var follows []string
	rows, err := db.c.Query("SELECT FOLLOWED_UUID FROM Follows WHERE UUID=?", uuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var followed string
		err = rows.Scan(&followed)
		if err != nil {
			return nil, err
		}
		follows = append(follows, followed)
	}
	if err != nil {
		return nil, err
	}
	return follows, nil
}
