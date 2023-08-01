package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Configurations struct {
	UserServiceName      string
	UserServiceDBUrl     string
	UserServicePort      string
	UserServiceSecretKey string
	ProxyBaseUrl         string

	TodoServiceName               string
	TodoServiceDBConnectionString string
	TodoServicePort               string

	TracingCollectorEndpoint string
}

func GetConfig(filepath string) *Configurations {
	err := godotenv.Load(filepath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	configurations := Configurations{
		ProxyBaseUrl: os.Getenv("PROXY_BASE_URL"),

		TracingCollectorEndpoint: os.Getenv("TRACING_COLLECTOR_ENDPOINT"),

		UserServiceName:      "users-service",
		UserServiceDBUrl:     os.Getenv("USER_SERVICE_DATABASE_URL"),
		UserServicePort:      os.Getenv("USER_SERVICE_PORT"),
		UserServiceSecretKey: os.Getenv("USER_SERVICE_SECRET"),

		TodoServiceName:               "todos-service",
		TodoServiceDBConnectionString: os.Getenv("TODO_SERVICE_CONNECTION_STRING"),
		TodoServicePort:               os.Getenv("TODO_SERVICE_PORT"),
	}

	return &configurations
}
