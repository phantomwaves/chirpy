package main

import (
	"fmt"
	"net/http"
)

type apiConfig struct {
	fileServerHits int
}

func (cfg *apiConfig) middlewareGetHits(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		cfg.fileServerHits++
		next.ServeHTTP(w, req)
	})
}

func (cfg *apiConfig) resetHitsHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	cfg.fileServerHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Counter reset"))
}

func (cfg *apiConfig) writeHitsHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %v", cfg.fileServerHits)))
}
func sendHealthResponse(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
