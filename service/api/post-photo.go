package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"log"
	"lucascutigliani.it/wasa/WasaPhoto/service/api/reqcontext"
	"lucascutigliani.it/wasa/WasaPhoto/service/utils"
	"net/http"
	"os"
)

func (rt *_router) postPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	_ = ps.ByName("")
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

	data, err := io.ReadAll(r.Body)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error reading request body")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Error reading request body"))
		return
	}
	err = r.Body.Close()
	if err != nil {
		ctx.Logger.WithError(err).Error("Error closing request body")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Error closing request body"))
		return
	}

	if len(data) < 100 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Photo too small"))
		return
	}

	if len(data) > 65000000 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Photo too big"))
		return
	}

	filePath := "./webui/uploads/"

	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(filePath, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	filePath = filePath + requestingUUID + "/"

	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(filePath, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	var extension string
	if isPNG(data) {
		extension = "png"
	} else if isJPEG(data) {
		extension = "jpg"
	} else {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid image format. Only PNG or JPEG are supported"))
	}

	photoUUID, err := rt.db.SetPhoto(requestingUUID, extension)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error setting photo")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Error setting photo"))
		return
	}

	filePath = filePath + photoUUID + "." + extension

	file, err := os.Create(filePath)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error creating file")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Error creating file"))
		return
	}

	_, err = file.Write(data)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error writing file")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Error writing file"))
		return
	}

	err = file.Close()
	if err != nil {
		ctx.Logger.WithError(err).Error("Error closing file")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Error closing file"))
		return
	}

	var responseMap = make(map[string]string)
	responseMap["photoID"] = photoUUID

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", "/uploads/"+requestingUUID+"/"+photoUUID+"."+extension)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(responseMap)
}

func isJPEG(data []byte) bool {
	return bytes.HasPrefix(data, []byte{0xFF, 0xD8, 0xFF})
}

func isPNG(data []byte) bool {
	return bytes.HasPrefix(data, []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A})
}
