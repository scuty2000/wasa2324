package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"lucascutigliani.it/wasa/WasaPhoto/service/api/reqcontext"
	"lucascutigliani.it/wasa/WasaPhoto/service/utils"
	"net/http"
	"strconv"
)

type PhotosResponse struct {
	UserPhotos      []string `json:"user-photos"`
	PaginationLimit int      `json:"paginationLimit"`
}

func (rt *_router) getUserPhotos(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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
		ctx.Logger.Warn(utils.InvalidBearer + requestingUUID)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Unauthorized: Authentication has failed."))
		return
	}

	requiredUUID := r.URL.Query().Get("userID")

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

	paginationIndex := r.URL.Query().Get("paginationIndex")
	if paginationIndex == "" {
		paginationIndex = "0"
	}

	paginationIndexInt, err := strconv.Atoi(paginationIndex)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad Request: paginationIndex must be an integer."))
		return
	}

	if paginationIndexInt < 0 {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad Request: paginationIndex must be a positive integer."))
		return
	}

	photos, photosCount, err := rt.db.GetUserPhotos(requiredUUID, paginationIndexInt)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error getting photos")
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("Internal Server Error: %s", err.Error())))
		return
	}

	paginationLimit := (photosCount - 1) / 10

	if len(photos) == 0 {
		photos = make([]string, 0)
	}

	var response PhotosResponse
	response.UserPhotos = photos
	response.PaginationLimit = paginationLimit

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
