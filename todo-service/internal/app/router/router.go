package router

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/riandyrn/otelchi"

	"github.com/go-chi/chi/v5"
	"github.com/olad5/productive-pulse/config"
	"github.com/olad5/productive-pulse/todo-service/internal/handlers"
)

func NewHttpRouter(todoHandler handlers.TodoHandler, configurations *config.Configurations) http.Handler {
	router := chi.NewRouter()
	router.Use(
		middleware.AllowContentType("application/json"),
		middleware.SetHeader("Content-Type", "application/json"),
	)
	router.Use(otelchi.Middleware(configurations.TodoServiceName, otelchi.WithChiRoutes(router)))

	router.Get("/todos/{id}", todoHandler.GetTodo)
	router.Get("/todos", todoHandler.GetTodos)
	router.Patch("/todos/{id}", todoHandler.UpdateTodo)
	router.Post("/todos", todoHandler.CreateTodo)
	return router
}
