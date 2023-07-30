package main

import (
	"context"
	"log"
	"net/http"

	"github.com/olad5/productive-pulse/config"
	"github.com/olad5/productive-pulse/pkg/app/server"
	"github.com/olad5/productive-pulse/todo-service/internal/app/router"
	"github.com/olad5/productive-pulse/todo-service/internal/handlers"
	"github.com/olad5/productive-pulse/todo-service/internal/infra/mongo"
	"github.com/olad5/productive-pulse/todo-service/internal/services/user"
	"github.com/olad5/productive-pulse/todo-service/internal/usecases/todos"
)

func main() {
	configurations := config.GetConfig(".env")
	ctx := context.Background()

	port := configurations.TodoServicePort
	todoRepo, err := mongo.NewMongoRepo(ctx, configurations.TodoServiceDBConnectionString)
	if err != nil {
		log.Fatal("Error Initializing todo Repo")
	}

	err = todoRepo.Ping(ctx)
	if err != nil {
		log.Fatal("Failed to ping todoRepo", err)
	}

	todoService, err := todos.NewTodoService(todoRepo, configurations)
	if err != nil {
		log.Fatal("Error Initializing TodoService")
	}

	userService, err := user.NewUserService(&http.Client{}, "http://localhost:"+configurations.UserServicePort)
	if err != nil {
		log.Fatal("Error Initializing UserService")
	}

	todoHandler, err := handlers.NewTodoHandler(*todoService, userService)
	if err != nil {
		log.Fatal("failed to create the TodoHandler: ", err)
	}

	appRouter := router.NewHttpRouter(*todoHandler)

	svr := server.CreateNewServer(appRouter)

	log.Printf("Server Listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, svr.Router))
}
