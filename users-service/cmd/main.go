package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/olad5/productive-pulse/config"
	"github.com/olad5/productive-pulse/users-service/internal/user"
)

func main() {
	configurations := config.GetConfig()
	ctx := context.Background()

	port := configurations.UserServicePort
	userRepo, err := user.NewPostgresRepo(ctx, configurations.UserServiceDBUrl)
	if err != nil {
		log.Fatal("Error Initializing User Repo")
	}

	err = userRepo.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("UserRepo initialized successfully")

	router := chi.NewRouter()

	log.Printf("Server Listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
