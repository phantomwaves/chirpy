package main

import (
	"log"
	"os"
	"strings"
)

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
