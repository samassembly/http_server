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

	//check for user matching email
	dbUser, err := cfg.databaseQueries.LoginUser(r.Context(), params.email)
	if err != nil {
		log.Printf("Incorrect Email or Password")
		w.WriteHeader(401)
		w.Write("Incorrect Email or Password")
		return
	}

	user := User{
		ID: dbUser.ID,
		CreatedAt: dbUser
	}
}