package main

import (
	"net/http"
	"log"
	"encoding/json"
	"github.com/samassembly/http_server/internal/auth"
	"github.com/samassembly/http_server/internal/database"
)

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {

	//decode request body
	type parameters struct {
        Email string `json:"email"`
		Password string `json:"password"`
    }

	decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
    }

	//hash password from params
	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error creating new user: %s", err)
		w.WriteHeader(500)
		return
	}

	dbParams := database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hash,
	}

	//call the database users function to create user with the decoded params.Email result
	user, err := cfg.databaseQueries.CreateUser(r.Context(), dbParams)
	if err != nil {
		log.Printf("Error creating new user: %s", err)
		w.WriteHeader(500)
		return
	}

	respBody := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}

    dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(dat)
	return
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	//decode request body
	type parameters struct {
        Email string `json:"email"`
		Password string `json:"password"`
    }
	

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Couldn't find token: %s", err)
		w.WriteHeader(500)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.servSecret)
	if err != nil {
		log.Printf("Couldn't get user for refresh token: %s", err)
		w.WriteHeader(401)
		return
	}

	decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err = decoder.Decode(&params)
    if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
    }

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error creating new user: %s", err)
		w.WriteHeader(500)
		return
	}

	//define params arguments
	update_params := database.UpdateUserParams{
		Email: params.Email,
		HashedPassword: hash,
		ID: userID,
	}

	user, err := cfg.databaseQueries.UpdateUser(r.Context(), update_params)
	if err != nil {
		log.Printf("Failed to update user: %s", err)
		w.WriteHeader(500)
		return
	}

	respBody := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
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