package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"lucascutigliani.it/wasa/WasaPhoto/service/api/reqcontext"
	"lucascutigliani.it/wasa/WasaPhoto/service/mocks"
	"lucascutigliani.it/wasa/WasaPhoto/service/utils"
	"net/http"
	"sort"
)

type QueryResult struct {
	Username string
	Uuid     string
	Distance int
}

func levenshtein(a, b string) int {
	if len(a) < len(b) {
		a, b = b, a
	}
	if len(b) == 0 {
		return len(a)
	}

	v0 := make([]int, len(b)+1)
	v1 := make([]int, len(b)+1)

	for i := range v0 {
		v0[i] = i
	}

	for i := 0; i < len(a); i++ {
		v1[0] = i + 1

		for j := 0; j < len(b); j++ {
			deletionCost := v0[j+1] + 1
			insertionCost := v1[j] + 1
			substitutionCost := v0[j]
			if a[i] != b[j] {
				substitutionCost++
			}

			// Trova il minimo tra deletionCost, insertionCost e substitutionCost
			minCost := deletionCost
			if insertionCost < minCost {
				minCost = insertionCost
			}
			if substitutionCost < minCost {
				minCost = substitutionCost
			}

			v1[j+1] = minCost
		}

		v0, v1 = v1, v0
	}

	return v0[len(b)]
}

func (rt *_router) getUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	_ = ps.ByName("")
	searchQuery := r.URL.Query().Get("searchQuery")
	if searchQuery == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad Request: searchQuery not provided."))
		return
	}

	if len(searchQuery) > 16 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad Request: searchQuery too long."))
		return
	}

	if len(searchQuery) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad Request: searchQuery too short."))
		return
	}

	matching, err := rt.db.SearchUsers(searchQuery)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error searching users")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal Server Error: " + err.Error()))
		return
	}

	var results []QueryResult
	for _, match := range matching {
		uuid := match[0]
		username := match[1]

		distance := levenshtein(username, searchQuery)
		results = append(results, QueryResult{Uuid: uuid, Username: username, Distance: distance})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Distance < results[j].Distance
	})

	users := make([]mocks.User, len(results))
	for i := 0; i < 20 && i < len(results); i++ {
		user, err := utils.MakeUserFromUUID(rt.db, results[i].Uuid)
		users[i] = *user
		if err != nil {
			ctx.Logger.WithError(err).Error("Error getting user")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("Internal Server Error: " + err.Error()))
			return
		}
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	var usersMap = make(map[string][]mocks.User)
	usersMap["users"] = users
	_ = json.NewEncoder(w).Encode(usersMap)

}
