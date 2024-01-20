package database

import (
	"database/sql"
	"errors"
)

func (db *appdbimpl) SearchUsers(searchQuery string) ([][]string, error) {
	rows, err := db.c.Query("SELECT UUID, USERNAME FROM Users WHERE USERNAME LIKE ?", "%"+searchQuery+"%")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	var results [][]string
	for rows.Next() {
		var uuid string
		var username string
		if err := rows.Scan(&uuid, &username); err != nil {
			return nil, err
		}
		result := []string{uuid, username}
		results = append(results, result)
	}
	if err != nil {
		return nil, err
	}

	return results, nil
}
