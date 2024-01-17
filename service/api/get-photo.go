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

func (rt *_router) getPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	requiredPhotoID := ps.ByName("photoID")
	ownerUUID, date, likes, comments, liked, err := rt.db.GetPhoto(requiredPhotoID, requestingUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("Not Found: Photo not found."))
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("Internal Server Error: %s", err.Error())))
		ctx.Logger.WithError(err).Error("Error retrieving photo")
		return
	}

	photo := &mocks.Photo{
		Uuid:          requiredPhotoID,
		Author:        ownerUUID,
		Date:          date,
		LikesCount:    likes,
		CommentsCount: comments,
		Liked:         liked,
	}

	hasPermission, err := utils.CheckUserAccess(rt.db, ctx, requestingUUID, ownerUUID)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(photo)
}
