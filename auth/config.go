package auth

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

func ReadConfig() Config {
	file, err := os.Open(".config.json")

	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	defer file.Close()

	bytes, err := io.ReadAll(file)

	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	var config Config

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}

	return config
}
