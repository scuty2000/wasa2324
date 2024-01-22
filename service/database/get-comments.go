package database

import (
	"lucascutigliani.it/wasa/WasaPhoto/service/mocks"
	"sort"
)

func (db *appdbimpl) GetPhotoComments(photoUUID string) ([]mocks.Comment, error) {
	rows, err := db.c.Query("SELECT COMMENT_UUID, OWNER_UUID, DATE, COMMENT_TEXT FROM Comments WHERE PHOTO_UUID =?", photoUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []mocks.Comment
	for rows.Next() {
		var commentUUID string
		var ownerUUID string
		var date string
		var commentText string
		if err := rows.Scan(&commentUUID, &ownerUUID, &date, &commentText); err != nil {
			return nil, err
		}
		result := mocks.Comment{
			Uuid:      commentUUID,
			OwnerUuid: ownerUUID,
			Date:      date,
			Text:      commentText,
		}
		results = append(results, result)
	}

	sort.SliceStable(results, func(i, j int) bool {
		return results[i].Date > results[j].Date
	})

	return results, err
}
