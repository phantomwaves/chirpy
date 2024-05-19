package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) newUserHandler(db *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		params, err := parseJsonRequest(req)
		if err != nil {
			log.Printf("Error parsing JSON request. User creation failed. %v", err)
			return
		}
		dbs, err := db.readDB()
		if err != nil {
			log.Println(err)
			http.Error(w, "Could not read database", http.StatusInternalServerError)
			return
		}
		user := dbs.CreateUser(params.Email)
		log.Printf("=============== %v, %v ==============", user.Id, user.Email)
		dbs.Users[user.Id] = user

		err2 := db.writeDB(dbs)

		if err2 != nil {
			log.Println(err2)
			http.Error(w, "Could not write user to database", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user) // send the created chirp in response
	}
}
