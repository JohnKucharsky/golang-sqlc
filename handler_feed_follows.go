package main

import (
	"encoding/json"
	"fmt"
	"github.com/JohnKucharsky/golang-sqlc/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (apiCfg *apiConfig) handlerCreateFeedFollow(
	w http.ResponseWriter,
	r *http.Request,
	user database.User,
) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id" validate:"required"`
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

	feedFollow, err := apiCfg.DB.CreateFeedFollow(
		r.Context(), database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			UserID:    user.ID,
			FeedID:    params.FeedID,
		},
	)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			fmt.Sprintf("Couldn't create feed follow:%v", err.Error()),
		)
		return
	}

	respondWithJSON(
		w, http.StatusCreated, databaseFeedFollowToFeedFollow(feedFollow),
	)
}

func (apiCfg *apiConfig) handlerGetFeedFollows(
	w http.ResponseWriter,
	r *http.Request,
	user database.User,
) {
	feedFollows, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			fmt.Sprintf("Couldn't get feed follows: %v", err.Error()),
		)
	}

	respondWithJSON(
		w,
		http.StatusCreated,
		databaseFeedFollowsToFeedFollows(feedFollows),
	)
}

func (apiCfg *apiConfig) handlerDeleteFeedFollow(
	w http.ResponseWriter,
	r *http.Request,
	user database.User,
) {
	feedFollowIdStr := chi.URLParam(r, "feedFollowId")
	feedFollowId, err := uuid.Parse(feedFollowIdStr)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			fmt.Sprintf("Couldn't delete feed follow:%v", err.Error()),
		)
		return
	}

	err = apiCfg.DB.DeleteFeedFollow(
		r.Context(),
		database.DeleteFeedFollowParams{ID: feedFollowId, UserID: user.ID},
	)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			fmt.Sprintf("Couldn't delete feed follow:%v", err.Error()),
		)
		return
	}

	respondWithJSON(
		w, http.StatusOK, struct {
		}{},
	)
}
