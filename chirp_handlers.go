package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type returnVals struct {
	Error       string `json:"error"`
	Valid       bool   `json:"valid"`
	CleanedBody string `json:"cleaned_body"`
}

func (cfg *apiConfig) validateChirp(w http.ResponseWriter, req *http.Request) {
	params, err := parseJsonRequest(req)
	if err != nil {
		log.Println("Error parsing request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	respBody := validateChirpLength(params)
	if !respBody.Valid {
		dat, err := writeJsonResponse(respBody)
		if err != nil {
			jsonMarshalError(w, err)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(dat)
	}
	respBody.CleanedBody = cleanChirp(params)

	dat, err := writeJsonResponse(respBody)
	if err != nil {
		jsonMarshalError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(dat)
	if err != nil {
		log.Println("error writing response body")
		return

	}
}

func (cfg *apiConfig) getChirpsHandler(db *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		chirps, err := cfg.getChirpsReq(db)
		if err != nil {
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(chirps)
	}

}

func (cfg *apiConfig) getChirpsReq(db *DB) ([]byte, error) {
	dbs, err := db.readDB()
	if err != nil {
		log.Println("Error reading chirps")
		return nil, err
	}
	cArray, err2 := dbs.GetChirps()
	if err2 != nil {
		return nil, err2
	}
	dat, err := json.Marshal(cArray)
	if err != nil {
		log.Println("Error marshalling chirps")
		return nil, err
	}
	return dat, nil
}

func (cfg *apiConfig) getChirpByIDHandler(db *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		chirpID := req.PathValue("ID")
		id, err := strconv.Atoi(chirpID)
		if err != nil {
			log.Println("Error converting chirpID to int")
			return
		}

		dbs, err2 := db.readDB()
		if err2 != nil {
			log.Println("Error reading DB")
			return
		}
		c, e := dbs.Chirps[id]
		if e == false {
			w.WriteHeader(404)
			log.Println("Chirp not found")
		}

		dat, err := json.Marshal(c)
		if err != nil {
			log.Println("Error marshalling chirp")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(dat)

	}
}

func (cfg *apiConfig) postChirpsHandler(db *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		params, err := parseJsonRequest(req)
		if err != nil {
			log.Println(err)
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		resp := validateChirpLength(params)
		resp.CleanedBody = cleanChirp(params)
		if !resp.Valid {
			http.Error(w, "Chirp is invalid", http.StatusBadRequest)
			return
		}

		dbs, err := db.readDB()
		if err != nil {
			log.Println(err)
			http.Error(w, "Could not read database", http.StatusInternalServerError)

		}
		chirp := dbs.CreateChirp(resp.CleanedBody)
		dbs.Chirps[chirp.Id] = chirp

		err2 := db.writeDB(dbs)

		if err2 != nil {
			log.Println(err2)
			http.Error(w, "Could not write to database", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(chirp) // send the created chirp in response
	}
}
