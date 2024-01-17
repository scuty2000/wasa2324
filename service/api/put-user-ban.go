package api

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"lucascutigliani.it/wasa/WasaPhoto/service/api/reqcontext"
	"lucascutigliani.it/wasa/WasaPhoto/service/utils"
	"net/http"
)

func (rt *_router) putUserBan(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	requestingUUID := ps.ByName("userID")
	bannedUserUUID := ps.ByName("bannedUserID")

	bearer := r.Header.Get("Authorization")
	if bearer == "" {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Unauthorized: Authentication not provided."))
		return
	}
	valid, err := utils.ValidateBearer(rt.db, ctx, requestingUUID, bearer)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error validating bearer token")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("Internal Server Error: %s", err.Error())))
		return
	}

	if !valid {
		ctx.Logger.Error("Authentication has failed")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Unauthorized: Authentication has failed."))
		return
	}

	if requestingUUID == bannedUserUUID {
		ctx.Logger.Error("Cannot ban yourself")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad Request: Cannot ban yourself."))
		return
	}

	err = rt.db.SetUserBan(requestingUUID, bannedUserUUID)
	if err != nil {
		if err.Error() != "UNIQUE constraint failed: Bans.UUID, Bans.BANNED_UUID" {
			ctx.Logger.WithError(err).Error("Error setting user ban")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprintf("Internal Server Error: %s", err.Error())))
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
