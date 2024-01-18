package api

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"lucascutigliani.it/wasa/WasaPhoto/service/api/reqcontext"
	"lucascutigliani.it/wasa/WasaPhoto/service/utils"
	"net/http"
)

func (rt *_router) deleteUserBan(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

		ctx.Logger.Warn("Invalid bearer token for user" + requestingUUID)
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Unauthorized: Authentication has failed."))
		return
	}

	if requestingUUID == bannedUserUUID {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad Request: userID and bannedUserID cannot be the same."))
		return
	}

	bans, err := rt.db.GetUserBans(requestingUUID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error getting user bans")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("Internal Server Error: %s", err.Error())))
		return
	}

	if !stringInSlice(bannedUserUUID, bans) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("Not Found: The required used is not banned"))
	}

	err = rt.db.DeleteUserBan(requestingUUID, bannedUserUUID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error deleting user ban")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("Internal Server Error: %s", err.Error())))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
