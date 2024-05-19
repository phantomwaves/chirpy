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

	db, _ := ensureDB()

	mux.Handle("/app/*", http.StripPrefix("/app/", apiCfg.middlewareGetHits(fs)))
	mux.HandleFunc("GET /api/healthz", sendHealthResponse)
	mux.HandleFunc("GET /admin/metrics", apiCfg.writeHitsHandler)
	mux.HandleFunc("GET /api/reset", apiCfg.resetHitsHandler)
	mux.HandleFunc("POST /api/chirps", apiCfg.postChirpsHandler(db))
	mux.HandleFunc("GET /api/chirps", apiCfg.getChirpsHandler(db))
	mux.HandleFunc("GET /api/chirps/{ID}", apiCfg.getChirpByIDHandler(db))

	server := &http.Server{
		Addr:    "localhost:" + port,
		Handler: mux}

	log.Printf("Server online at http://%s/app\n", server.Addr)
	log.Fatal(server.ListenAndServe())

}
