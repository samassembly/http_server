package main

import (
	"net/http"
	"log"
	"encoding/json"
	"github.com/samassembly/http_server/internal/auth"
	//"github.com/samassembly/http_server/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
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

	//check for user matching email
	dbUser, err := cfg.databaseQueries.LoginUser(r.Context(), params.Email)
	if err != nil {
		log.Printf("Incorrect Email or Password")
		w.WriteHeader(401)
		w.Write([]byte("Incorrect Email or Password"))
		return
	}

	//compare input password to stored hash
	ok := auth.CheckPasswordHash(dbUser.HashedPassword, params.Password)
	if ok != nil {
		log.Printf("Passwords do not match: %s", err)
		w.WriteHeader(401)
		return
	}

	//craft response
	respBody := User{
		ID: dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email: dbUser.Email,
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