package main

import (
	"net/http"
	"log"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.cfgPlatform != "dev" {
		w.WriteHeader(403)
		w.Write([]byte("Access forbidden"))
		w.Write([]byte(cfg.cfgPlatform))
	}

	_, err := cfg.databaseQueries.ResetUsers(r.Context())
	if err != nil {
		log.Printf("Error resetting users: %s", err)
		w.WriteHeader(500)
		return
	}


	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}