package main

import (
	"context"
	"log"
	"net/http"

	"github.com/olad5/productive-pulse/config"
	"github.com/olad5/productive-pulse/pkg/app/server"
	"github.com/olad5/productive-pulse/pkg/utils"
	"github.com/olad5/productive-pulse/todo-service/internal/app/router"
	"github.com/olad5/productive-pulse/todo-service/internal/handlers"
	"github.com/olad5/productive-pulse/todo-service/internal/infra/mongo"

	"github.com/olad5/productive-pulse/todo-service/internal/services/user"
	"github.com/olad5/productive-pulse/todo-service/internal/usecases/todos"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	configurations := config.GetConfig(".env")
	ctx := context.Background()

	tracerProvider, err := utils.NewTracerProvider(configurations.TodoServiceName, configurations.TracingCollectorEndpoint)
	if err != nil {
		log.Fatal("JaegerTraceProvider failed to Initialize", err)
	}

	mongoMonitor := otelmongo.NewMonitor(otelmongo.WithTracerProvider(tracerProvider))
	tracer := tracerProvider.Tracer(configurations.TodoServiceName)

	todoRepo, err := mongo.NewMongoRepo(ctx, mongoMonitor, configurations.TodoServiceDBConnectionString)
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

	userServiceClient := &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	userService, err := user.NewUserService(userServiceClient, configurations.ProxyBaseUrl)
	if err != nil {
		log.Fatal("Error Initializing UserService")
	}

	todoHandler, err := handlers.NewTodoHandler(*todoService, userService, tracer)
	if err != nil {
		log.Fatal("failed to create the TodoHandler: ", err)
	}

	appRouter := router.NewHttpRouter(*todoHandler, configurations)

	svr := server.CreateNewServer(appRouter)

	port := configurations.TodoServicePort

	log.Printf("Server Listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, svr.Router))
}
