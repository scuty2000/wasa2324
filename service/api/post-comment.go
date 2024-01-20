package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"lucascutigliani.it/wasa/WasaPhoto/service/api/reqcontext"
	"lucascutigliani.it/wasa/WasaPhoto/service/utils"
	"net/http"
	"regexp"
)

func (rt *_router) postComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			ctx.Logger.WithError(err).Error("Error closing request body")
		}
	}(r.Body)

	var jsonMap = make(map[string]string)
	content, err := io.ReadAll(r.Body)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error reading request body")
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Error reading request body")
		return
	}
	if !json.Valid(content) || len(content) == 0 {
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Invalid JSON")
		return
	}
	err = json.Unmarshal(content, &jsonMap)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error decoding session login json")
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Error decoding session login json")
		return
	}
	keys := make([]string, 0, len(jsonMap))
	for k := range jsonMap {
		keys = append(keys, k)
	}
	if len(keys) != 2 {
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("JSON not conforming to schema")
		return
	}

	var hasAllFiels = true
	for _, key := range keys {
		if key != "text" && key != "issuer" {
			hasAllFiels = false
		}
	}

	if !hasAllFiels {
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("JSON not conforming to schema")
		return
	}

	issuerUUID := jsonMap["issuer"]
	if issuerUUID == "" {
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Bad request: issuer is empty")
		return
	}

	pattern := "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"

	re, err := regexp.Compile(pattern)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error compiling regular expression")
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode("An error occurred while processing your request")
		return
	}

	if !re.MatchString(issuerUUID) {
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Issuer " + issuerUUID + " is not a valid UUID")
		return
	}

	bearer := r.Header.Get("Authorization")
	if bearer == "" {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Unauthorized: Authentication not provided."))
		return
	}

	valid, err := utils.ValidateBearer(rt.db, ctx, issuerUUID, bearer)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error validating bearer token")
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("Internal Server Error: %s", err.Error())))
		return
	}

	if !valid {
		ctx.Logger.Warn(utils.InvalidBearer + issuerUUID)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Unauthorized: Authentication has failed."))
		return
	}

	photoUUID := ps.ByName("photoID")
	photoOwnerUUID, _, _, _, _, _, err := rt.db.GetPhoto(photoUUID, issuerUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.Header().Set("content-type", "text/plain")
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode("Not Found: Photo not found")
			return
		} else {
			ctx.Logger.WithError(err).Error("Error retrieving photo")
			w.Header().Set("content-type", "text/plain")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode("Internal Server Error: An error occurred while processing your request")
			return
		}
	}

	hasPermission, err := utils.CheckUserAccess(rt.db, ctx, issuerUUID, photoOwnerUUID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error checking user access")
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode("Internal Server Error: An error occurred while processing your request")
		return
	}

	if !hasPermission {
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode("Forbidden: You don't have permission to access this resource")
		return
	}

	commentUUID, date, err := rt.db.SetComment(issuerUUID, photoUUID, jsonMap["text"])
	if err != nil {
		ctx.Logger.WithError(err).Error("Error setting comment")
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode("Internal Server Error: An error occurred while processing your request")
		return
	}

	var returnJson = make(map[string]string)
	returnJson["id"] = commentUUID
	returnJson["date"] = date
	returnJson["text"] = jsonMap["text"]

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(returnJson)
}
