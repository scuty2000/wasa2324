package utils

import (
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/mocks"
)

func MakeUserFromUUID(db database.AppDatabase, uuid string) (*mocks.User, error) {
	username, err := db.GetUserByUUID(uuid)
	if err != nil {
		return nil, err
	}

	// TODO IMPLEMENTS MISSING FIELDS RETRIEVAL

	return &mocks.User{
		Uuid:           uuid,
		Username:       username,
		FollowersCount: 0,
		FollowingCount: 0,
		PhotosCount:    0,
		IsBanned:       false,
	}, err
}
