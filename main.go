package main

import (
	"database/sql"
	"net/http"
	"sync/atomic"
	"os"
	"time"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/samassembly/http_server/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	databaseQueries *database.Queries
	cfgPlatform string
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type Chirp struct {
	ID		  uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body	  string	`json:"body"`
	User_ID   uuid.UUID `json:"user_id"`
}

func main() {
	//load .env into vars
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")

	//open connection to database
	db, _ := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)

	// Create new apiConfig
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		databaseQueries: dbQueries,
		cfgPlatform: platform,
	}

	// Create new ServeMux
	mux := http.NewServeMux()

	//handler functions 
	//app
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	//api
	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{id}", apiCfg.handlerGetChirp)

	mux.HandleFunc("POST /api/users", apiCfg.handlerUsers)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirps)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	//admin
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	// Create server struct
	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	//Start server
	server.ListenAndServe()
}
