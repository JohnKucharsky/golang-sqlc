package main

import (
	"encoding/json"
	"fmt"
	"github.com/JohnKucharsky/golang-sqlc/internal/database"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (apiCfg *apiConfig) handlerCreateFeed(
	w http.ResponseWriter,
	r *http.Request,
	user database.User,
) {
	type parameters struct {
		Name string `json:"name" validate:"required"`
		URL  string `json:"url" validate:"required"`
	}
	params := parameters{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			fmt.Sprintf("Error parsing JSON:%v", err),
		)
		return
	}

	if err := validate.Struct(params); err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			fmt.Sprintf("Error validating JSON:%v", err.Error()),
		)
		return

	}

	feed, err := apiCfg.DB.CreateFeed(
		r.Context(), database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name:      params.Name,
			Url:       params.URL,
			UserID:    user.ID,
		},
	)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			fmt.Sprintf("Couldn't create feed:%v", err),
		)
		return
	}

	respondWithJSON(
		w, http.StatusCreated, databaseFeedToFeed(feed),
	)
}

func (apiCfg *apiConfig) handlerGetFeeds(
	w http.ResponseWriter,
	r *http.Request,
) {
	feeds, err := apiCfg.DB.GetFeeds(r.Context())

	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			fmt.Sprintf("Couldn't get feeds:%v", err),
		)
		return
	}

	respondWithJSON(w, http.StatusOK, databaseFeedsToFeeds(feeds))
}
