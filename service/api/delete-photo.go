package api

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"lucascutigliani.it/wasa/WasaPhoto/service/api/reqcontext"
	"lucascutigliani.it/wasa/WasaPhoto/service/utils"
	"net/http"
	"os"
)

func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	photoUUID := ps.ByName("photoID")
	requestingUUID := r.Header.Get("X-Requesting-User-UUID")

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
		ctx.Logger.Warn(utils.InvalidBearer + requestingUUID)
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Unauthorized: Authentication has failed."))
		return
	}

	ownerUUID, _, extension, _, _, _, err := rt.db.GetPhoto(photoUUID, requestingUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("Not Found: Photo not found."))
			return
		}
		ctx.Logger.WithError(err).Error("Error getting photo")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("Internal Server Error: %s", err.Error())))
		return
	}

	if requestingUUID != ownerUUID {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("Forbidden: You are not the owner of this photo."))
		return
	}

	err = rt.db.DeletePhoto(photoUUID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error deleting photo")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("Internal Server Error: %s", err.Error())))
		return
	}

	filePath := "./uploads/" + ownerUUID + "/" + photoUUID + "." + extension
	err = os.Remove(filePath)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error deleting photo file")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("Internal Server Error: %s", err.Error())))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
