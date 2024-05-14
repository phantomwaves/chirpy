package main

import (
	"log"
	"net/http"
)

const port string = "8080"

func main() {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("."))
	apiCfg := apiConfig{}

	mux.Handle("/app/*", http.StripPrefix("/app/", apiCfg.middlewareGetHits(fs)))
	mux.HandleFunc("/api/healthz", sendHealthResponse)
	mux.HandleFunc("/admin/metrics", apiCfg.writeHitsHandler)
	mux.HandleFunc("/api/reset", apiCfg.resetHitsHandler)

	server := &http.Server{
		Addr:    "localhost:" + port,
		Handler: mux}

	log.Printf("Server online at http://%s/app\n", server.Addr)
	log.Fatal(server.ListenAndServe())

}
