package database

import "fmt"

func (db *appdbimpl) GetUserPhotos(uuid string, offsetMultiplier int) ([]string, int, error) {
	var photos []string

	var photoCount int
	err := db.c.QueryRow(fmt.Sprintf("SELECT COUNT(1) FROM Photos WHERE OWNER_UUID = '%s'", uuid)).Scan(&photoCount)
	if err != nil {
		return nil, 0, err
	}

	var formattedQuery string
	if offsetMultiplier == 0 {
		formattedQuery = fmt.Sprintf("SELECT UUID FROM Photos WHERE OWNER_UUID = '%s' ORDER BY date DESC LIMIT 10", uuid)
	} else {
		formattedQuery = fmt.Sprintf("SELECT UUID FROM Photos WHERE OWNER_UUID = '%s' ORDER BY date DESC LIMIT 10 OFFSET %d", uuid, offsetMultiplier*10)
	}

	rows, err := db.c.Query(formattedQuery)
	if err != nil {
		return nil, 0, err
	}
	if rows.Err() != nil {
		return nil, 0, rows.Err()
	}
	defer rows.Close()
	for rows.Next() {
		var follower string
		err = rows.Scan(&follower)
		if err != nil {
			return nil, 0, err
		}
		photos = append(photos, follower)
	}
	if err != nil {
		return nil, 0, err
	}
	return photos, photoCount, nil
}
