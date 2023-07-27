package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Configurations struct {
	UserServiceDBUrl     string
	UserServicePort      string
	UserServiceSecretKey string
}

func GetConfig(filepath string) *Configurations {
	err := godotenv.Load(filepath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	configurations := Configurations{
		UserServiceDBUrl:     os.Getenv("USER_SERVICE_DATABASE_URL"),
		UserServicePort:      os.Getenv("USER_SERVICE_PORT"),
		UserServiceSecretKey: os.Getenv("USER_SERVICE_SECRET"),
	}

	return &configurations
}
