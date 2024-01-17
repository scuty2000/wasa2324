package api

import (
	"database/sql"
	"errors"
	"fmt"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/utils"
	"github.com/julienschmidt/httprouter"
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
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Unauthorized: Authentication has failed."))
		return
	}

	ownerUUID, _, _, _, _, err := rt.db.GetPhoto(photoUUID, requestingUUID)
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

	filePath := "./webui/uploads/" + ownerUUID + "/" + photoUUID + ".jpg"
	err = os.Remove(filePath)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error deleting photo")
		if errors.Is(err, os.ErrNotExist) {
			filePath = "./webui/uploads/" + ownerUUID + "/" + photoUUID + ".png"
			err = os.Remove(filePath)
			if err != nil {
				ctx.Logger.WithError(err).Error("Error deleting photo1")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(fmt.Sprintf("Internal Server Error: %s", err.Error())))
				return
			}
		} else {
			ctx.Logger.WithError(err).Error("Error deleting photo2")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprintf("Internal Server Error: %s", err.Error())))
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
