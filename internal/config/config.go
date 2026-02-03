package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ConfigStr struct {
	DatabaseURL string
	Port        string
}

func LoadEnvironmentalVariables() (*ConfigStr, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
		return nil, err
	}

	dbURL := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")

	configStr := &ConfigStr{
		DatabaseURL: dbURL,
		Port:        port,
	}
	return configStr, nil
}