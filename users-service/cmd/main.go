package main

import (
	"context"
	"log"
	"net/http"

	"github.com/olad5/productive-pulse/config"
	"github.com/olad5/productive-pulse/users-service/internal/app/router"
	"github.com/olad5/productive-pulse/users-service/internal/app/server"
	"github.com/olad5/productive-pulse/users-service/internal/handlers"
	"github.com/olad5/productive-pulse/users-service/internal/infra/postgres"
	"github.com/olad5/productive-pulse/users-service/internal/usecases/users"
)

func main() {
	configurations := config.GetConfig(".env")
	ctx := context.Background()

	port := configurations.UserServicePort
	userRepo, err := postgres.NewPostgresRepo(ctx, configurations.UserServiceDBUrl)
	if err != nil {
		log.Fatal("Error Initializing User Repo")
	}

	err = userRepo.Ping(ctx)
	if err != nil {
		log.Fatal("Failed to ping UserRepo", err)
	}

	userService, err := users.NewUserService(userRepo, configurations)
	if err != nil {
		log.Fatal("Error Initializing UserService")
	}

	userHandler, err := handlers.NewHandler(*userService)
	if err != nil {
		log.Fatal("failed to create the User handler: ", err)
	}

	appRouter := router.NewHttpRouter(*userHandler)

	svr := server.CreateNewServer(appRouter)

	log.Printf("Server Listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, svr.Router))
}
