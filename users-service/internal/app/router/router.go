package router

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/riandyrn/otelchi"

	"github.com/go-chi/chi/v5"
	"github.com/olad5/productive-pulse/config"
	"github.com/olad5/productive-pulse/users-service/internal/handlers"
)

func NewHttpRouter(userHandler handlers.UserHandler, configurations *config.Configurations) http.Handler {
	router := chi.NewRouter()
	router.Use(
		middleware.AllowContentType("application/json"),
		middleware.SetHeader("Content-Type", "application/json"),
	)
	router.Use(otelchi.Middleware(configurations.UserServiceName, otelchi.WithChiRoutes(router)))

	router.Get("/users/auth", userHandler.Auth)
	router.Post("/users/login", userHandler.Login)
	router.Post("/users", userHandler.Register)
	return router
}
