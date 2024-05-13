package main

import (
	"fmt"
	"log"
	"net/http"
)

type apiConfig struct {
	fileServerHits int
}

func (cfg *apiConfig) middlewareGetHits(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileServerHits++
		next.ServeHTTP(w, req)
	})
}

func (cfg *apiConfig) resetHitsHandler(w http.ResponseWriter, req *http.Request) {
	cfg.fileServerHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Counter reset"))
}

func (cfg *apiConfig) writeHitsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(fmt.Sprintf("Hits: %v", cfg.fileServerHits)))
}

func send_response(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	const port = "8080"
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("."))
	apiCfg := apiConfig{}

	mux.Handle("/app/", http.StripPrefix("/app/", apiCfg.middlewareGetHits(fs)))
	mux.HandleFunc("/healthz", send_response)
	mux.HandleFunc("/metrics", apiCfg.writeHitsHandler)
	mux.HandleFunc("/reset", apiCfg.resetHitsHandler)

	server := &http.Server{
		Addr:    "localhost:" + port,
		Handler: mux}

	log.Printf("Server online at http://%s/app\n", server.Addr)
	log.Fatal(server.ListenAndServe())

}
