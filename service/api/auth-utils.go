package api

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
)

func AuthUser(rt *_router, ctx reqcontext.RequestContext, username string) (string, string, bool, error) {
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
	uuid, err = rt.db.GetUserByName(username)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			ctx.Logger.WithError(err).Error("Error getting user UUID")
			return "", "", false, err
		}
		uuid, err = rt.db.CreateUser(username)
		if err != nil {
			ctx.Logger.WithError(err).Error("Error creating user")
			return "", "", false, err
		}
		created = true
		ctx.Logger.Info("Created user ", username)
	}

	return token, uuid, created, nil
}

func ValidateBearer(rt *_router, ctx reqcontext.RequestContext, uuid string, bearer string) (bool, error) {
	var token string
	var err error
	token, err = rt.db.GetUserSession(uuid)
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
