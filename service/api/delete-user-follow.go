package api

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"lucascutigliani.it/wasa/WasaPhoto/service/api/reqcontext"
	"lucascutigliani.it/wasa/WasaPhoto/service/utils"
	"net/http"
)

func (rt *_router) deleteUserFollow(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	requestingUUID := ps.ByName("userID")
	followedUUID := ps.ByName("followedID")

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

	if requestingUUID == followedUUID {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad Request: Cannot unfollow yourself."))
		return
	}

	follows, err := rt.db.GetUserFollows(requestingUUID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error getting user follows")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("Internal Server Error: %s", err.Error())))
		return
	}

	if !stringInSlice(followedUUID, follows) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("Not Found: The required used is not followed"))
	}

	err = rt.db.DeleteUserFollow(requestingUUID, followedUUID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error deleting user follow")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("Internal Server Error: %s", err.Error())))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
