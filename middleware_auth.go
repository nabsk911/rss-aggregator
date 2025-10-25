package main

import (
	"fmt"
	"net/http"

	"github.com/nabsk911/rss-aggregator/internal/auth"
	"github.com/nabsk911/rss-aggregator/internal/db"
)

type authedHandler func(http.ResponseWriter, *http.Request, db.User)

func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		apiKey, err := auth.GetAPIKey(r.Header)

		if err != nil {
			writeJSON(w, http.StatusUnauthorized, envelope{"error": fmt.Sprint(err)})
			return
		}

		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)

		if err != nil {
			writeJSON(w, http.StatusBadRequest, envelope{"error": fmt.Sprint(err)})
			return
		}

		handler(w, r, user)
	}
}
