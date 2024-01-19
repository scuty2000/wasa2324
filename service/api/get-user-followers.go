package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"lucascutigliani.it/wasa/WasaPhoto/service/api/reqcontext"
	"lucascutigliani.it/wasa/WasaPhoto/service/mocks"
	"lucascutigliani.it/wasa/WasaPhoto/service/utils"
	"net/http"
)

func (rt *_router) getUserFollowers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	requestingUUID := r.Header.Get("X-Requesting-User-UUID")
	if requestingUUID == "" {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad Request: requesting userID not provided in header."))
		return
	}

	bearer := r.Header.Get("Authorization")
	if bearer == "" {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Unauthorized: Authentication not provided."))
		return
	}
	valid, err := utils.ValidateBearer(rt.db, ctx, requestingUUID, bearer)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error validating bearer token")
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("Internal Server Error: %s", err.Error())))
		return
	}

	if !valid {
		ctx.Logger.Warn("Invalid bearer token for user" + requestingUUID)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Unauthorized: Authentication has failed."))
		return
	}

	requiredUUID := ps.ByName("userID")

	_, err = rt.db.GetUserByUUID(requiredUUID)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("Not Found: User not found."))
			return
		}
	}

	hasPermission, err := utils.CheckUserAccess(rt.db, ctx, requestingUUID, requiredUUID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error checking user access")
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("Internal Server Error: %s", err.Error())))
		return
	}

	if !hasPermission {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("Forbidden: Insufficient permissions to perform this action. You may have tried to perform an operation on another user data or tried to retrieve data from a user that banned you."))
		return
	}

	var followersMap = make(map[string][]mocks.User)

	followers, err := rt.db.GetUserFollowers(requiredUUID)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			followersMap["followers"] = []mocks.User{}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(followersMap)
			return
		}
		ctx.Logger.WithError(err).Error("Error getting user followers")
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal Server Error: Error getting user followers."))
		return
	}

	for _, follower := range followers {
		followerUser, err := utils.MakeUserFromUUID(rt.db, ctx, follower, requestingUUID)
		if err != nil {
			ctx.Logger.WithError(err).Error("Error getting user followers")
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("Internal Server Error: Error getting user followers."))
			return
		}
		followersMap["followers"] = append(followersMap["followers"], *followerUser)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(followersMap)
}
