package utils

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"lucascutigliani.it/wasa/WasaPhoto/service/api/reqcontext"
	"lucascutigliani.it/wasa/WasaPhoto/service/database"
)

func AuthUser(db database.AppDatabase, ctx reqcontext.RequestContext, username string) (string, string, bool, error) {
	bytes := make([]byte, 128)
	_, err := rand.Read(bytes)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error generating bearer token")
		return "", "", false, err
	}
	token := "Bearer " + hex.EncodeToString(bytes)
	ctx.Logger.Info("Generated bearer token for user ", username)

	var uuid string
	var created = false
	uuid, err = db.GetUserByName(username)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			ctx.Logger.WithError(err).Error("Error getting user UUID")
			return "", "", false, err
		}
		uuid, err = db.CreateUser(username)
		if err != nil {
			ctx.Logger.WithError(err).Error("Error creating user")
			return "", "", false, err
		}
		created = true
		ctx.Logger.Info("Created user ", username)
	}

	err = db.SetSession(uuid, token)
	if err != nil {
		return "", "", false, err
	}

	return token, uuid, created, nil
}

func ValidateBearer(db database.AppDatabase, ctx reqcontext.RequestContext, uuid string, bearer string) (bool, error) {
	var token string
	var err error
	token, err = db.GetUserSession(uuid)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			ctx.Logger.WithError(err).Error("Error getting bearer token")
			return false, err
		}
		ctx.Logger.WithError(errors.New("bearer token not found")).Error("Error getting bearer token")
		return false, nil
	}
	if token != bearer {
		ctx.Logger.WithError(errors.New("bearer token not matching")).Error("Bearer token not matching")
		return false, nil
	}
	return true, nil
}

func CheckUserAccess(db database.AppDatabase, ctx reqcontext.RequestContext, requestingUUID string, requestedUUID string) (bool, error) {
	banned, err := db.GetUserBans(requestedUUID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error getting bans")
		return false, err
	}
	for _, ban := range banned {
		if ban == requestingUUID {
			return false, nil
		}
	}
	return true, nil
}
