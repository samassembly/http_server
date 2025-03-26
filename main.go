package main

import (
	"net/http"
)



func main() {
	// Create new ServeMux
	mux := http.NewServeMux()

	// Create new Fileserver
	fileServer := http.FileServer(http.Dir("."))

	//Handler functions here (blank for first step)
	mux.Handle("/", fileServer)
	mux.Handle("/assets/logo.png", fileServer)

	// Create server struct
	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	//Start server
	server.ListenAndServe()
}

