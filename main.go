package main

import (
	"flag"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

const port string = "8080"

func main() {
	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if *dbg == true {
		log.Println("Debug mode enabled -- removing database.json")
		err := os.Remove("database.json")
		if err != nil {
			log.Fatal("Error removing existing database", err)
		}
	}
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("."))
	apiCfg := apiConfig{
		jwtSecret: os.Getenv("JWT_KEY"),
	}

	db, _ := ensureDB()

	mux.Handle("/app/*", http.StripPrefix("/app/", apiCfg.middlewareGetHits(fs)))
	mux.HandleFunc("GET /api/healthz", sendHealthResponse)
	mux.HandleFunc("GET /admin/metrics", apiCfg.writeHitsHandler)
	mux.HandleFunc("GET /api/reset", apiCfg.resetHitsHandler)
	mux.HandleFunc("POST /api/chirps", apiCfg.postChirpsHandler(db))
	mux.HandleFunc("GET /api/chirps", apiCfg.getChirpsHandler(db))
	mux.HandleFunc("GET /api/chirps/{ID}", apiCfg.getChirpByIDHandler(db))
	mux.HandleFunc("POST /api/users", apiCfg.newUserHandler(db))
	mux.HandleFunc("PUT /api/users", apiCfg.updateUserHandler(db))
	mux.HandleFunc("POST /api/login", apiCfg.loginHandler(db))

	server := &http.Server{
		Addr:    "localhost:" + port,
		Handler: mux}

	log.Printf("Server online at http://%s/app\n", server.Addr)
	log.Fatal(server.ListenAndServe())

}
