package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

type apiConfig struct {
	fileServerHits int
}

func getHTMLFromFile(filepath string, replace map[string]string) []byte {
	contents, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println("Error: file not found")
		return []byte("")
	}
	if len(replace) == 0 {
		return contents
	}
	for k, v := range replace {
		if strings.Contains(string(contents), k) {
			contents = []byte(strings.ReplaceAll(string(contents), k, v))
		}
	}
	return contents

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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(getHTMLFromFile("metrics.html",
		map[string]string{"{XX}": fmt.Sprint(cfg.fileServerHits)}))

	// w.Write([]byte(fmt.Sprintf("Hits: %v", cfg.fileServerHits)))
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
