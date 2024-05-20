package main

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type UserResponse struct {
	Email string `json:"email"`
	Id    int    `json:"id"`
}

func (cfg *apiConfig) newUserHandler(db *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		params, err := parseJsonRequest(req)
		if err != nil {
			log.Printf("Error parsing JSON request. User creation failed. %v", err)
			http.Error(w, "JSON parsing error", http.StatusInternalServerError)
			return
		}
		if params.Password == "" {
			log.Printf("User creation failed. Password is empty.")
			http.Error(w, "User creation failed. Password is empty.", http.StatusInternalServerError)
			return
		}
		dbs, err := db.readDB()
		if err != nil {
			log.Printf("Error reading DB. User creation failed. %v", err)
			http.Error(w, "DB read error", http.StatusInternalServerError)
			return
		}
		_, err3 := dbs.checkEmail(params.Email)
		if err3 == nil {
			log.Printf("User with that email already exists.")
			http.Error(w, "User with that email already exists.", http.StatusInternalServerError)
			return
		}

		user := dbs.CreateUser(params.Email, params.Password)
		dbs.Users[user.Id] = user

		err2 := db.writeDB(dbs)

		if err2 != nil {
			log.Println(err2)
			http.Error(w, "Could not write user to database", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		noPW := UserResponse{
			Id:    user.Id,
			Email: params.Email,
		}
		x, _ := json.Marshal(noPW)
		w.Write(x)
	}
}

func (cfg *apiConfig) loginHandler(db *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		params, err := parseJsonRequest(req)
		if err != nil {
			log.Printf("Error parsing JSON request. User creation failed. %v", err)
			http.Error(w, "JSON parsing error", http.StatusInternalServerError)
			return
		}
		if params.Password == "" {
			log.Printf("Login failed. Password is empty.")
			http.Error(w, "Login failed. Password is empty.", http.StatusInternalServerError)
			return
		}
		dbs, err := db.readDB()
		if err != nil {
			log.Printf("Error reading DB. Login failed. %v", err)
			http.Error(w, "DB Error. Login failed.", http.StatusInternalServerError)
			return
		}
		id, err := dbs.checkEmail(params.Email)
		if err != nil {
			log.Printf("User not found. %v", err)
			http.Error(w, "User not found.", http.StatusInternalServerError)
		}
		hp := dbs.Users[id].Password
		log.Printf("id: %v, email: %v, hashed password: %v.", id, params.Email, hp)
		valid := bcrypt.CompareHashAndPassword([]byte(hp), []byte(params.Password))
		if valid != nil {
			log.Println("Incorrect password")
			http.Error(w, "Incorrect password", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		noPW := UserResponse{
			Id:    id,
			Email: dbs.Users[id].Email,
		}
		x, _ := json.Marshal(noPW)
		w.Write(x)

	}
}
