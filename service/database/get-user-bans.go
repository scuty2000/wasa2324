package database

func (db *appdbimpl) GetUserBans(uuid string) ([]string, error) {
	var bans []string
	rows, err := db.c.Query("SELECT BANNED_UUID FROM Bans WHERE UUID=?", uuid)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	defer rows.Close()
	for rows.Next() {
		var ban string
		err = rows.Scan(&ban)
		if err != nil {
			return nil, err
		}
		bans = append(bans, ban)
	}
	return bans, nil
}
