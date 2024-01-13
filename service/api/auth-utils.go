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
	}

	return token, uuid, created, nil
}
