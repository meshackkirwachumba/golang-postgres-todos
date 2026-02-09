package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ConfigStr struct {
	DatabaseURL string
	Port        string
	JWTSecret   string
}

func LoadEnvironmentalVariables() (*ConfigStr, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
		return nil, err
	}

	dbURL := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")
	jwtSecret := os.Getenv("JWT_SECRET")

	configStr := &ConfigStr{
		DatabaseURL: dbURL,
		Port:        port,
		JWTSecret:   jwtSecret,
	}
	return configStr, nil
}