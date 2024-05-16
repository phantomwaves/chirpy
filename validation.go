package main

import (
	"log"
	"strings"
)

func validateChirpLength(params parameters) returnVals {
	if len(params.Body) > 140 {
		log.Printf("Chirp is too long.")
		return returnVals{
			Error: "Chirp is too long",
			Valid: false,
		}
	}
	return returnVals{Valid: true}
}

func cleanChirp(params parameters) string {
	var cleaned string
	for _, word := range strings.Split(params.Body, " ") {
		wLower := strings.ToLower(word)
		if wLower == "kerfuffle" || wLower == "sharbert" || wLower == "fornax" {
			cleaned += " " + "****"
		} else {
			cleaned += " " + word
		}
	}
	cleaned = strings.TrimLeft(cleaned, " ")
	return cleaned
}
