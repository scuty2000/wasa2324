package api

import (
	"encoding/json"
	"fmt"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/utils"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) getUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	requestingUUID := r.Header.Get("X-Requesting-User-UUID")
	if requestingUUID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad Request: requesting userID not provided in header."))
		return
	}

	bearer := r.Header.Get("Authorization")
	if bearer == "" {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Unauthorized: Authentication not provided."))
		return
	}
	valid, err := utils.ValidateBearer(rt.db, ctx, requestingUUID, bearer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("Internal Server Error: %s", err.Error())))
		ctx.Logger.WithError(err).Error("Error validating bearer token")
		return
	}

	if !valid {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Unauthorized: Authentication has failed."))
		ctx.Logger.Error("Authentication has failed")
		return
	}

	requiredUUID := ps.ByName("userID")
	hasPermission, err := utils.CheckUserAccess(rt.db, ctx, requestingUUID, requiredUUID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("Internal Server Error: %s", err.Error())))
		ctx.Logger.WithError(err).Error("Error checking user access")
		return
	}

	if !hasPermission {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("Forbidden: Insufficient permissions to perform this action. You may have tried to perform an operation on another user data or tried to retrieve data from a user that banned you."))
		return
	}

	user, err := utils.MakeUserFromUUID(rt.db, requiredUUID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("Internal Server Error: %s", err.Error())))
		ctx.Logger.WithError(err).Error("Error getting user")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user)
	return
}
