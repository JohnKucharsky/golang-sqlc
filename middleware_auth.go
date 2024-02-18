package main

import (
	"fmt"
	"github.com/JohnKucharsky/golang-sqlc/internal/auth"
	"github.com/JohnKucharsky/golang-sqlc/internal/database"
	"net/http"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		handler(w, r, user)
	}
}
