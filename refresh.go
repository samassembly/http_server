package main

import (
	"net/http"
	"time"
	"log"
	"encoding/json"
	"github.com/samassembly/http_server/internal/auth"
)

// TokenResponse defines the structure of the JSON response
type TokenResponse struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Couldn't find token: %s", err)
		w.WriteHeader(500)
		return
	}

	user, err := cfg.databaseQueries.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		log.Printf("Couldn't get user for refresh token: %s", err)
		w.WriteHeader(401)
		return
	}

	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.servSecret,
		time.Hour,
	)
	if err != nil {
		log.Printf("Couldn't validate token: %s", err)
		w.WriteHeader(500)
		return
	}

	//craft response
	respBody := TokenResponse{
		Token: accessToken,
	}

	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
	return
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Couldn't find token: %s", err)
		w.WriteHeader(500)
		return
	}

	_, err = cfg.databaseQueries.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		log.Printf("Couldn't revoke session: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
