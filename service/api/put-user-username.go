package api

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"lucascutigliani.it/wasa/WasaPhoto/service/api/reqcontext"
	"lucascutigliani.it/wasa/WasaPhoto/service/utils"
	"net/http"
	"regexp"
)

func (rt *_router) putUserUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	requiredUUID := ps.ByName("userID")

	bearer := r.Header.Get("Authorization")
	if bearer == "" {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Unauthorized: Authentication not provided."))
		return
	}
	valid, err := utils.ValidateBearer(rt.db, ctx, requiredUUID, bearer)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error validating bearer token")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("Internal Server Error: %s", err.Error())))
		return
	}

	if !valid {
		ctx.Logger.Warn("Invalid bearer token for user" + requiredUUID)
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("Unauthorized: Authentication has failed."))
		return
	}

	var jsonMap = make(map[string]string)
	content, err := io.ReadAll(r.Body)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error reading request body")
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Error reading request body"))
		return
	}
	if !json.Valid(content) || len(content) == 0 {
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid JSON"))
		return
	}
	err = json.Unmarshal(content, &jsonMap)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error decoding session login json")
		w.Header().Set("content-type", "text/plain")
		_, _ = w.Write([]byte("Error decoding session login json"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	keys := make([]string, 0, len(jsonMap))
	for k := range jsonMap {
		keys = append(keys, k)
	}
	if len(keys) != 1 || keys[0] != "username" {
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("JSON not conforming to schema"))
		return
	}

	username := jsonMap["username"]
	if username == "" {
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Username is empty"))
		return
	}

	pattern := "^[a-zA-Z0-9_]*$"

	re, err := regexp.Compile(pattern)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error compiling regular expression")
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("An error occurred while processing your request"))
		return
	}

	if !re.MatchString(username) {
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf("Name \"%s\" does not match pattern", username)))
		return
	}

	if len(username) < 3 {
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf("Name \"%s\" is too short", username)))
		return
	}

	if len(username) > 16 {
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf("Name \"%s\" is too long", username)))
		return
	}

	err = rt.db.UpdateUsername(requiredUUID, username)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error updating username")
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("An error occurred while processing your request"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
