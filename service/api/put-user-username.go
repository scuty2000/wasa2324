package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/utils"
	"github.com/julienschmidt/httprouter"
	"io"
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
		ctx.Logger.Error("Authentication has failed")
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
		ctx.Logger.WithError(errors.New("invalid JSON string")).Error("Invalid JSON")
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
		ctx.Logger.WithError(errors.New("json not conforming to schema")).Error("JSON not conforming to schema")
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("JSON not conforming to schema"))
		return
	}

	username := jsonMap["username"]
	if username == "" {
		ctx.Logger.WithError(errors.New("username is empty")).Error("Username is empty")
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
		ctx.Logger.WithError(errors.New("username not matching regex")).Error("Name does not match pattern")
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf("Name \"%s\" does not match pattern", username)))
		return
	}

	if len(username) < 3 {
		ctx.Logger.WithError(errors.New("username too short")).Error("Name too short")
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf("Name \"%s\" is too short", username)))
		return
	}

	if len(username) > 16 {
		ctx.Logger.WithError(errors.New("username too long")).Error("Name too long")
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
