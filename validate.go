package main

import (
	"net/http"
	"log"
	"encoding/json"
	"regexp"
)

func handlerValidate(w http.ResponseWriter, r *http.Request) {

	//profane const
    profane := "****"


	//decode request body
	type parameters struct {
        Body string `json:"body"`
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
	
	type returnVals struct {
		Cleaned_body string `json:"cleaned_body"`
    }
	respBody := returnVals{
		Cleaned_body: cleanedBody,
	}

    dat, err := json.Marshal(respBody)
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
