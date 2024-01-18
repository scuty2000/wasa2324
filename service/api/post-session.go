package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	"lucascutigliani.it/wasa/WasaPhoto/service/api/reqcontext"
	"lucascutigliani.it/wasa/WasaPhoto/service/utils"
	"net/http"
	"regexp"
)

func (rt *_router) postSession(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	_ = ps.ByName("")
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			ctx.Logger.Error(err)
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
	if len(keys) != 1 || keys[0] != "name" {
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("JSON not conforming to schema")
		return
	}

	username := jsonMap["name"]
	if username == "" {
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Empty username")
		return
	}

	pattern := "^[a-zA-Z0-9_]*$"

	re, err := regexp.Compile(pattern)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error compiling regular expression")
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode("An error occurred while processing your request")
		return
	}

	if !re.MatchString(username) {
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Name " + username + " does not match pattern")
		return
	}

	if len(username) < 3 {
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Name too short")
		return
	}

	if len(username) > 16 {
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("Name too long")
		return
	}

	bearer, uuid, created, err := utils.AuthUser(rt.db, ctx, username)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error generating bearer token")
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode("An error occurred while processing your request")
		return
	}

	authMap := make(map[string]string)
	authMap["identifier"] = uuid
	authMap["token"] = bearer

	if created {
		ctx.Logger.Info("User " + username + " created (\"" + uuid + "\")")
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(authMap)
		return
	}
	ctx.Logger.Info("User " + username + " logged in")
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(authMap)
}
