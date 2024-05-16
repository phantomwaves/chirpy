package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func parseJsonRequest(req *http.Request) (parameters, error) {
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		return params, err
	}

	return params, nil
}

func writeJsonResponse(resp returnVals) ([]byte, error) {
	dat, err := json.Marshal(resp)
	if err != nil {

		return nil, err
	}
	return dat, nil
}

func jsonMarshalError(w http.ResponseWriter, err error) {
	log.Printf("Error marshalling JSON: %s", err)
	w.WriteHeader(http.StatusInternalServerError)
	respBody := returnVals{
		Error: "Something went wrong",
		Valid: false,
	}
	dat, _ := json.Marshal(respBody)
	_, _ = w.Write(dat)
}
