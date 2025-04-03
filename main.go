package main

import (
	"database/sql"
	"net/http"
	"sync/atomic"
	"os"
	"github.com/joho/godotenv"
	"github.com/samassembly/http_server/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	databaseQueries *database.Queries
}

func main() {
	//load .env into vars
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	//open connection to database
	db, _ := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)

	// Create new apiConfig
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		databaseQueries: dbQueries,
	}

	// Create new ServeMux
	mux := http.NewServeMux()

	//handler functions 
	//app
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	//api
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidate)
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
