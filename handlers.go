package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type apiConfig struct {
	fileServerHits int
}

type parameters struct {
	Body string `json:"body"`
}

type returnVals struct {
	Error string `json:"error"`
	Valid bool   `json:"valid"`
}

func getHTMLFromFile(filepath string, replace map[string]string) ([]byte, error) {
	contents, err := os.ReadFile(filepath)
	if err != nil {
		log.Print("file not found")
		return nil, err
	}
	if len(replace) == 0 {
		return contents, nil
	}
	for k, v := range replace {
		if strings.Contains(string(contents), k) {
			contents = []byte(strings.ReplaceAll(string(contents), k, v))
		}
	}
	return contents, nil

}

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

func (cfg *apiConfig) validateChirp(w http.ResponseWriter, req *http.Request) {
	params, err := parseJsonRequest(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(params.Body) > 140 {
		log.Printf("Chirp is too long.")

		respBodyTooLong := returnVals{
			Error: "Chirp is too long",
			Valid: false,
		}

		dat, err := writeJsonResponse(respBodyTooLong)
		if err != nil {
			jsonMarshalError(w, err)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(dat)
		if err != nil {
			return
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")

	respBodyValid := returnVals{
		Valid: true,
	}

	dat, err := writeJsonResponse(respBodyValid)
	if err != nil {
		jsonMarshalError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(dat)

	if err != nil {
		return

	}
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
