package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nabsk911/rss-aggregator/internal/db"
)

func (apiCfg *apiConfig) handleCreateFeed(w http.ResponseWriter, r *http.Request, user db.User) {

	type createFeedRequest struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	req := createFeedRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, envelope{"error": "Invalid feed payload!"})
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), db.CreateFeedParams{
		ID:        uuid.New(),
		Name:      req.Name,
		Url:       req.URL,
		UserID:    user.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, envelope{"error": "Internal server error!"})
		return
	}
	writeJSON(w, http.StatusCreated, envelope{"data": feed, "message": "Feed created successfully!"})

}

func (apiCfg *apiConfig) handleGetFeeds(w http.ResponseWriter, r *http.Request) {

	feeds, err := apiCfg.DB.GetFeeds(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, envelope{"error": "Internal server error!"})
		return
	}
	writeJSON(w, http.StatusOK, envelope{"data": feeds})
}
