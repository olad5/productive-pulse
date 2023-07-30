package router

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/olad5/productive-pulse/todo-service/internal/handlers"
)

func NewHttpRouter(todoHandler handlers.TodoHandler) http.Handler {
	router := chi.NewRouter()
	router.Use(
		middleware.AllowContentType("application/json"),
		middleware.SetHeader("Content-Type", "application/json"),
	)

	router.Get("/todos/{id}", todoHandler.GetTodo)
	router.Get("/todos", todoHandler.GetTodos)
	router.Patch("/todos/{id}", todoHandler.UpdateTodo)
	router.Post("/todos", todoHandler.CreateTodo)
	return router
}
