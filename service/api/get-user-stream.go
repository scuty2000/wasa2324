package api

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"lucascutigliani.it/wasa/WasaPhoto/service/api/reqcontext"
	"lucascutigliani.it/wasa/WasaPhoto/service/mocks"
	"lucascutigliani.it/wasa/WasaPhoto/service/utils"
	"net/http"
	"strconv"
)

type StreamResponse struct {
	Photos          []mocks.Photo `json:"photos"`
	PaginationLimit int           `json:"paginationLimit"`
}

func (rt *_router) getUserStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	userID := ps.ByName("userID")

	bearer := r.Header.Get("Authorization")
	if bearer == "" {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Unauthorized: Authentication not provided."))
		return
	}
	valid, err := utils.ValidateBearer(rt.db, ctx, userID, bearer)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error validating bearer token")
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("Internal Server Error: %s", err.Error())))
		return
	}

	if !valid {
		ctx.Logger.Warn(utils.InvalidBearer + userID)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Unauthorized: Authentication has failed."))
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

	photos, photosCount, err := rt.db.GetPaginatedPhotos(userID, paginationIndexInt)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error getting photos")
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("Internal Server Error: %s", err.Error())))
		return
	}

	paginationLimit := (photosCount - 1) / 10

	if len(photos) == 0 {
		photos = []mocks.Photo{}
	}

	var response StreamResponse
	response.Photos = photos
	response.PaginationLimit = paginationLimit

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
