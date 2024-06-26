package main

import (
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"sync"
)

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type parameters struct {
	Body             string `json:"body"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	Token            string `json:"token"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

type DB struct {
	path string
	mux  *sync.RWMutex
}

func (dbs *DBStructure) CreateChirp(body string) Chirp {
	nextID := dbs.nextID("chirp")
	return Chirp{Body: body, Id: nextID}
}
func (dbs *DBStructure) CreateUser(email, password string) User {
	nextID := dbs.nextID("user")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error generating password hash: %v", err)
	}
	return User{Email: email, Password: string(hashedPassword), Id: nextID}
}

func (dbs *DBStructure) GetChirps() ([]Chirp, error) {
	chirpsList := make([]Chirp, 0, len(dbs.Chirps))
	for _, c := range dbs.Chirps {
		chirpsList = append(chirpsList, c)
	}
	return chirpsList, nil
}

func newDB(path string) (*DB, error) {
	initialData :=
		`{
    "chirps": {},
    "users": {}
	}` // Initial JSON structure for an empty 'chirps' map
	err := os.WriteFile(path, []byte(initialData), 0777)
	if err != nil {
		return nil, errors.New("could not create DB")
	}
	db := &DB{path: path, mux: &sync.RWMutex{}}
	return db, nil
}

// EnsureDB return nil if database.json already exists, else create and return it
func ensureDB() (*DB, error) {
	_, err := os.ReadFile("database.json")
	path := "database.json"
	if err != nil {
		log.Println("database does not exist. creating new database...")
		db, err2 := newDB(path)
		if err2 != nil {
			return nil, err2
		}
		return db, nil
	}
	log.Println("database exists. reading database...")
	db := &DB{path: path, mux: &sync.RWMutex{}}
	return db, nil
}

func (db *DB) writeDB(dbs DBStructure) error {
	c, err := json.Marshal(dbs)
	if err != nil {
		return errors.New("could not encode chirps")
	}
	db.mux.Lock()
	defer db.mux.Unlock()
	err = os.WriteFile(db.path, c, 0777)
	if err != nil {
		return errors.New("could not write chirps to db")
	}
	log.Println("wrote chirps to db")

	return nil
}

func (db *DB) readDB() (DBStructure, error) {
	if db.path == "" {
		log.Println("DB path is empty!")
		return DBStructure{}, errors.New("db path is empty")
	}
	data, err := os.ReadFile(db.path)
	if err != nil {
		log.Printf("Error reading file: %v\n", err)
		return DBStructure{}, errors.New("could not read chirps from db")
	}
	var dbs DBStructure
	err = json.Unmarshal(data, &dbs)
	if err != nil {
		log.Printf("Error unmarshalling chirps: %v\n", err)
		return DBStructure{}, errors.New("could not decode chirps from db")
	}
	return dbs, nil
}

func (dbs *DBStructure) nextID(x string) int {
	var nextID int
	switch x {
	case "chirp":
		if len(dbs.Chirps) > 0 {
			for id := range dbs.Chirps {
				if id >= nextID {
					nextID = id + 1
				}
			}
		} else {
			nextID = 1
		}
	case "user":
		if len(dbs.Users) > 0 {
			for id := range dbs.Users {
				if id >= nextID {
					nextID = id + 1
				}
			}
		} else {
			nextID = 1
		}

	}

	return nextID
}

func (dbs *DBStructure) checkEmail(email string) (int, error) {
	if email == "" {
		return 0, errors.New("invalid email")
	}
	for _, u := range dbs.Users {
		if u.Email == email {
			return u.Id, nil
		}
	}
	return 0, errors.New("email not in use")
}

func (dbs *DBStructure) updateUser(email, password string, id int) {
	hp, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error generating password hash: %v", err)
	}
	u := User{
		Id:       id,
		Email:    email,
		Password: string(hp),
	}
	dbs.Users[id] = u
}
