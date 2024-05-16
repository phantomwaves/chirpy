package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestPostChirp(t *testing.T) {
	// Ensure fresh database state
	if err := os.Remove("database.json"); err != nil && !os.IsNotExist(err) {
		t.Fatalf("Could not remove database: %v", err)
	}

	db, err := ensureDB()
	if err != nil {
		t.Fatalf("Could not initialize database: %v", err)
	}

	createChirp(t, "I had something interesting for breakfast", db)
	createChirp(t, "What about second breakfast?", db)
	createChirp(t, "Supper? Dinner? Do you know about those?", db)
}

func createChirp(t *testing.T, body string, db *DB) {
	chirp := map[string]string{"body": body}
	jsonData, err := json.Marshal(chirp)
	if err != nil {
		t.Fatalf("Error marshalling json: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/chirps", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	cfg := apiConfig{}
	handler := cfg.postChirpsHandler(db)
	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusCreated {
		t.Errorf("Expected status code %v, got %v", http.StatusCreated, status)
	}

	// Further checks on response body
	var respBody map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("Error unmarshalling response: %v", err)
	}

	if respBody["body"] != body {
		t.Errorf("Expected body %v, got %v", body, respBody["body"])
	}

	if _, ok := respBody["id"].(float64); !ok {
		t.Errorf("Expected id to be a number")
	}
}
