package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/nabsk911/rss-aggregator/internal/db"
)

func (apiCfg *apiConfig) handleCreateFeedFollows(w http.ResponseWriter, r *http.Request, user db.User) {
	type createFeedFollowRequest struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	req := createFeedFollowRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, envelope{"error": "Invalid feed follow payload!"})
		return
	}

	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), db.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    req.FeedID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, envelope{"error": "Internal server error!"})
		return
	}
	writeJSON(w, http.StatusCreated, envelope{"data": feedFollow, "message": "Feed follow created successfully!"})

}

func (apiCfg *apiConfig) handleGetFeedFollows(w http.ResponseWriter, r *http.Request, user db.User) {

	feedFollows, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, envelope{"error": "Internal server error!"})
		return
	}
	writeJSON(w, http.StatusOK, envelope{"data": feedFollows})
}

func (apiCfg *apiConfig) handleDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user db.User) {
	chiURLParam := chi.URLParam(r, "feedFollowID")
	feedID, err := uuid.Parse(chiURLParam)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, envelope{"error": "Invalid feed id!"})
		return
	}
	err = apiCfg.DB.DeleteFeedFollow(r.Context(), db.DeleteFeedFollowParams{
		ID:     feedID,
		UserID: user.ID,
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, envelope{"error": "Internal server error!"})
		return
	}
	writeJSON(w, http.StatusNoContent, envelope{"message": "Feed follow deleted successfully!"})
}
