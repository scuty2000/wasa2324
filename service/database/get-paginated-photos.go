package database

import (
	"database/sql"
	"errors"
	"lucascutigliani.it/wasa/WasaPhoto/service/mocks"
	"sort"
	"strings"
)

func (db *appdbimpl) GetPaginatedPhotos(requestingUUID string, offsetMultiplier int) ([]mocks.Photo, int, error) {

	var bannedFrom []string
	rows, err := db.c.Query("SELECT UUID from Bans WHERE BANNED_UUID = ?", requestingUUID)
	if err != nil {
		return nil, 0, err
	}
	for rows.Next() {
		var bannedUUID string
		if err := rows.Scan(&bannedUUID); err != nil {
			return nil, 0, err
		}
		bannedFrom = append(bannedFrom, bannedUUID)
	}
	bannedFrom = append(bannedFrom, requestingUUID)

	var following []string
	rows, err = db.c.Query("SELECT FOLLOWED_UUID from Follows WHERE UUID = ?", requestingUUID)
	if err != nil {
		return nil, 0, err
	}
	for rows.Next() {
		var followedUUID string
		if err := rows.Scan(&followedUUID); err != nil {
			return nil, 0, err
		}
		following = append(following, followedUUID)
	}

	uuidsForSQL := "'" + strings.Join(bannedFrom, "', '") + "'"
	followingForSQL := "'" + strings.Join(following, "', '") + "'"

	var photosCount int
	err = db.c.QueryRow("SELECT COUNT(*) FROM Photos WHERE OWNER_UUID NOT IN (?) AND OWNER_UUID IN (?)", uuidsForSQL, followingForSQL).Scan(&photosCount)

	if offsetMultiplier == 0 {
		rows, err = db.c.Query("SELECT * FROM Photos WHERE OWNER_UUID NOT IN (?) AND OWNER_UUID IN (?) ORDER BY date DESC LIMIT 10", uuidsForSQL, followingForSQL)
	} else {
		rows, err = db.c.Query("SELECT * FROM Photos WHERE OWNER_UUID NOT IN (?) AND OWNER_UUID IN (?) ORDER BY date DESC LIMIT 10 OFFSET ?", uuidsForSQL, followingForSQL, offsetMultiplier*10)
	}

	if err != nil {
		return nil, 0, err
	}

	var results []mocks.Photo
	for rows.Next() {
		var uuid string
		var ownerUUID string
		var date string
		var extension string
		err := rows.Scan(&uuid, &ownerUUID, &date, &extension)
		if err != nil {
			return nil, 0, err
		}
		var likes int
		err = db.c.QueryRow("SELECT COUNT(*) FROM Likes WHERE PHOTO_UUID=?", uuid).Scan(&likes)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				likes = 0
			} else {
				return nil, 0, err
			}
		}
		var comments int
		err = db.c.QueryRow("SELECT COUNT(*) FROM Comments WHERE PHOTO_UUID=?", uuid).Scan(&comments)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				comments = 0
			} else {
				return nil, 0, err
			}
		}
		var isLiked int
		err = db.c.QueryRow("SELECT COUNT(*) FROM Likes WHERE PHOTO_UUID=? AND USER_UUID=?", uuid, requestingUUID).Scan(&isLiked)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				isLiked = 0
			} else {
				return nil, 0, err
			}
		}
		result := mocks.Photo{
			Uuid:          uuid,
			Author:        ownerUUID,
			Date:          date,
			Extension:     extension,
			LikesCount:    likes,
			CommentsCount: comments,
			Liked:         isLiked == 1,
		}
		results = append(results, result)
	}

	sort.SliceStable(results, func(i, j int) bool {
		return results[i].Date > results[j].Date
	})

	return results, photosCount, err
}
