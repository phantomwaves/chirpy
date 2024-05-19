package main

import (
	"fmt"
	"log"
	"net/http"
)

func (cfg *apiConfig) middlewareGetHits(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileServerHits++
		next.ServeHTTP(w, req)
	})
}

func (cfg *apiConfig) resetHitsHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits = 0
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Counter reset"))
	if err != nil {
		return
	}
	return
}

func (cfg *apiConfig) writeHitsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	body, err := getHTMLFromFile("metrics.html", map[string]string{"{XX}": fmt.Sprint(cfg.fileServerHits)})
	if err != nil {
		log.Println("error reading html file")
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(body)
	if err != nil {
		log.Println("error writing response body")
		return
	}
	return
}

func sendHealthResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		log.Println("error writing response body")
		return
	}
	return
}
