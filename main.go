package main

import (
	// "fmt"
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    "localhost:" + port,
		Handler: mux}

	log.Printf("Server online at http://%s\n", server.Addr)
	log.Fatal(server.ListenAndServe())

}
