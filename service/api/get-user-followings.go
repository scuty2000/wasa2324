package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/mocks"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/utils"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) getUserFollowings(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("Internal Server Error: %s", err.Error())))
		ctx.Logger.WithError(err).Error("Error validating bearer token")
		return
	}

	if !valid {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Unauthorized: Authentication has failed."))
		ctx.Logger.Error("Authentication has failed")
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
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("Internal Server Error: %s", err.Error())))
		ctx.Logger.WithError(err).Error("Error checking user access")
		return
	}

	if !hasPermission {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("Forbidden: Insufficient permissions to perform this action. You may have tried to perform an operation on another user data or tried to retrieve data from a user that banned you."))
		return
	}

	var followingsMap = make(map[string][]mocks.User)

	followings, err := rt.db.GetUserFollows(requiredUUID)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			followingsMap["followings"] = []mocks.User{}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(followingsMap)
			return
		}
		ctx.Logger.WithError(err).Error("Error getting user followers")
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal Server Error: Error getting user followers."))
		return
	}

	for _, following := range followings {
		followingUser, err := utils.MakeUserFromUUID(rt.db, following)
		if err != nil {
			ctx.Logger.WithError(err).Error("Error getting user followers")
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("Internal Server Error: Error getting user followers."))
			return
		}
		followingsMap["followings"] = append(followingsMap["followings"], *followingUser)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(followingsMap)
}
