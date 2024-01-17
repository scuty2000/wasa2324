package api

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"lucascutigliani.it/wasa/WasaPhoto/service/api/reqcontext"
	"lucascutigliani.it/wasa/WasaPhoto/service/utils"
	"net/http"
)

func (rt *_router) putUserLike(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	requestingUUID := ps.ByName("userID")
	photoUUID := ps.ByName("photoID")

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

	ownerUUID, _, _, _, _, _, err := rt.db.GetPhoto(photoUUID, requestingUUID)
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

	hasPermission, err := utils.CheckUserAccess(rt.db, ctx, requestingUUID, ownerUUID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error checking user access")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("Internal Server Error: %s", err.Error())))
		return
	}

	if !hasPermission {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("Forbidden: Forbidden: Insufficient permissions to perform this action. You may have tried to perform an operation on another user data or tried to retrieve data from a user that banned you."))
		return
	}

	err = rt.db.SetUserLike(requestingUUID, photoUUID)
	if err != nil {
		if err.Error() != "UNIQUE constraint failed: Likes.UUID, Likes.PHOTO_UUID" {
			ctx.Logger.WithError(err).Error("Error setting user like")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprintf("Internal Server Error: %s", err.Error())))
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
