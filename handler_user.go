package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nabsk911/rss-aggregator/internal/db"
)

func (apiCfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type userRequest struct {
		Name string `json:"name"`
	}

	userReq := userRequest{}

	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		log.Printf("Decoding user data: %v", err)
		writeJSON(w, http.StatusBadRequest, envelope{"error": "Invalid user payload!"})
		return
	}
	createdUser, err := apiCfg.DB.CreateUser(r.Context(), db.CreateUserParams{
		ID:        uuid.New(),
		Name:      userReq.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		log.Printf("Creating user: %v", err)
		writeJSON(w, http.StatusInternalServerError, envelope{"error": "Internal server error!"})
		return
	}
	writeJSON(w, http.StatusCreated, envelope{"data": createdUser, "message": "User created successfully!"})
}
func (apiCfg *apiConfig) handleGetUser(w http.ResponseWriter, r *http.Request, user db.User) {

	writeJSON(w, http.StatusAccepted, envelope{"data": user})
}

func (apiCfg *apiConfig) handleGetPostsForUser(w http.ResponseWriter, r *http.Request, user db.User) {

	posts, err := apiCfg.DB.GetPostsForUser(r.Context(), db.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, envelope{"error": "Internal server error!"})
		return
	}
	writeJSON(w, http.StatusOK, envelope{"data": posts})
}
