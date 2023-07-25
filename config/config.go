package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Configurations struct {
	UserServiceDBUrl string
	UserServicePort  string
}

func GetConfig() *Configurations {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	configurations := Configurations{
		UserServiceDBUrl: os.Getenv("USER_SERVICE_DATABASE_URL"),
		UserServicePort:  os.Getenv("USER_SERVICE_PORT"),
	}

	return &configurations
}
