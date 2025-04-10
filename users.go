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
