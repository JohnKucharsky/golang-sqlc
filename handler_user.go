package main

import (
	"encoding/json"
	"fmt"
	"github.com/JohnKucharsky/golang-sqlc/internal/auth"
	"github.com/JohnKucharsky/golang-sqlc/internal/database"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (apiCfg *apiConfig) handlerCreateUser(
	w http.ResponseWriter,
	r *http.Request,
) {
	type parameters struct {
		Name string `json:"name" validate:"required"`
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

	user, err := apiCfg.DB.CreateUser(
		r.Context(), database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name:      params.Name,
		},
	)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			fmt.Sprintf("Couldn't crate user:%v", err),
		)
		return
	}

	respondWithJSON(
		w, http.StatusCreated, databaseUserToUser(user),
	)
}

func (apiCfg *apiConfig) handlerGetUser(
	w http.ResponseWriter,
	r *http.Request,
) {
	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		respondWithError(
			w,
			http.StatusForbidden,
			fmt.Sprintf("Auth error: %v", err.Error()),
		)
		return
	}

	user, err := apiCfg.DB.GetUserByApiKey(r.Context(), apiKey)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			fmt.Sprintf("Couldn't get the user: %v", err.Error()),
		)
		return
	}

	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}
