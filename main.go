package main

import (
	"net/http"
)



func main() {
	// Create new ServeMux
	mux := http.NewServeMux()

	//Handler functions here (blank for first step)

	// Create server struct
	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	//Start server
	server.ListenAndServe()
}

