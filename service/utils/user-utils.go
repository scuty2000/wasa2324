package utils

import (
	"lucascutigliani.it/wasa/WasaPhoto/service/api/reqcontext"
	"lucascutigliani.it/wasa/WasaPhoto/service/database"
	"lucascutigliani.it/wasa/WasaPhoto/service/mocks"
)

func MakeUserFromUUID(db database.AppDatabase, ctx reqcontext.RequestContext, uuid string, requestingUUID string) (*mocks.User, error) {
	username, err := db.GetUserByUUID(uuid)
	if err != nil {
		return nil, err
	}

	followersCount, err := db.GetFollowersCount(uuid)
	if err != nil {
		return nil, err
	}

	followingCount, err := db.GetFollowingCount(uuid)
	if err != nil {
		return nil, err
	}

	photosCount, err := db.GetUserPhotosCount(uuid)
	if err != nil {
		return nil, err
	}

	var isBanned bool
	if requestingUUID == "" {
		isBanned = false
	} else {
		isBanned, err = CheckUserAccess(db, ctx, uuid, requestingUUID)
		if err != nil {
			return nil, err
		}
	}

	return &mocks.User{
		Uuid:           uuid,
		Username:       username,
		FollowersCount: followersCount,
		FollowingCount: followingCount,
		PhotosCount:    photosCount,
		IsBanned:       isBanned,
	}, err
}
