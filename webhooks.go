package main

import (
	"net/http"
	"log"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/samassembly/http_server/internal/auth"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	
	type WebhookData struct {
		UserID uuid.UUID `json:"user_id"`
		}

	type WebhookRequest struct {
		Event string `json:"event"`
		Data WebhookData `json:"data"`
	}
		
apiKey, err := auth.GetAPIKey(r.Header)
if err != nil {
	log.Printf("Error retrieving API Key: %s", err)
	w.WriteHeader(500)
	return
}

if apiKey != cfg.polkaKey {
	log.Printf("Invalid API Key")
	w.WriteHeader(401)
	return
}

decoder := json.NewDecoder(r.Body)
reqBody := WebhookRequest{}
err = decoder.Decode(&reqBody)
if err != nil {
	log.Printf("Error decoding webhook request: %s", err)
	w.WriteHeader(500)
	return
}

if reqBody.Event != "user.upgraded" {
	log.Printf("Invalid Event Recieved: %s", reqBody.Event)
	w.WriteHeader(204)
	return
}

_, err = cfg.databaseQueries.UpgradeUser(r.Context(), reqBody.Data.UserID)
if err != nil {
	log.Printf("Error upgrading user: %s", err)
	w.WriteHeader(404)
	return
}

log.Printf("User Upgraded Successfully")
w.WriteHeader(204)
return
}