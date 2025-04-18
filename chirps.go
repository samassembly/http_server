package main

import (
	"net/http"
	"log"
	"encoding/json"
	"regexp"
	"github.com/samassembly/http_server/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {

	//profane const
    profane := "****"

	//decode request body
	type parameters struct {
        Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
    }

	decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
    }

	//encode response body
	//invalid response length
	if len(params.Body) > 140 {
		type returnVals struct {
			Error string `json:"error"`
		}
		respBody := returnVals{
			Error: "Chirp is too long",
		}
		dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
		return
	}

	//valid response length
	pattern := `(?i)\b(kerfuffle|sharbert|fornax)\b`
    re := regexp.MustCompile(pattern)
    cleanedBody := re.ReplaceAllString(params.Body, profane)

	chirpParams := database.CreateChirpParams{
		Body: cleanedBody,
		UserID: params.UserID,
	}

	//call the database chirp function to create chirp with the clean body
	chirp, err := cfg.databaseQueries.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		log.Printf("Error creating new chirp: %s", err)
		w.WriteHeader(500)
		return
	}
	
	respBody := Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		User_ID: chirp.UserID,
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

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.databaseQueries.GetChirps(r.Context())
	if err != nil {
		log.Printf("Error getting chirps: %s", err)
		w.WriteHeader(500)
		return
	}

	chirps := []Chirp{}	
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			User_ID:    dbChirp.UserID,
			Body:      dbChirp.Body,
		})
	}

	dat, err := json.Marshal(chirps)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
	return
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	id_str := r.PathValue("id")
	id, err := uuid.Parse(id_str)
	if err != nil {
		log.Printf("Error parsing chirp id: %s", err)
	 	w.WriteHeader(500)
	 	return
	}

	dbChirp, err := cfg.databaseQueries.GetChirp(r.Context(), id)
	if err != nil {
		log.Printf("Error getting chirp: %s", err)
		w.WriteHeader(404)
		return
	}

	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		User_ID:   dbChirp.UserID,
		Body:      dbChirp.Body,
	}	
	

	dat, err := json.Marshal(chirp)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
	return
}