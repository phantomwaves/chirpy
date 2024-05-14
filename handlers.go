package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type apiConfig struct {
	fileServerHits int
}

type returnVals struct {
	Error string `json:"error"`
	Valid bool   `json:"valid"`
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

}

func (cfg *apiConfig) validateChirp(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(params.Body) > 140 {
		log.Printf("Chirp is too long.")
		respBody := returnVals{
			Error: "Chirp is too long",
			Valid: false,
		}
		dat, err := json.Marshal(respBody)
		if err != nil {
			jsonMarshalError(w, err)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(dat)
	}

	respBody := returnVals{
		Valid: true,
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		jsonMarshalError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(dat)

}

func jsonMarshalError(w http.ResponseWriter, err error) {
	log.Printf("Error marshalling JSON: %s", err)
	w.WriteHeader(http.StatusInternalServerError)
	respBody := returnVals{
		Error: "Something went wrong",
		Valid: false,
	}
	dat, _ := json.Marshal(respBody)
	w.Write(dat)
	return
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
