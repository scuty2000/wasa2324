package database

func (db *appdbimpl) GetUserBans(uuid string) ([]string, error) {
	var bans []string
	rows, err := db.c.Query("SELECT BANNED_UUID FROM Bans WHERE UUID=?", uuid)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var ban string
		err = rows.Scan(&ban)
		if err != nil {
			return nil, err
		}
		bans = append(bans, ban)
	}
	err = rows.Close()
	if err != nil {
		return nil, err
	}
	return bans, nil
}
