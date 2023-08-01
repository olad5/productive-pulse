package handlers

import (
	"errors"

	"github.com/olad5/productive-pulse/users-service/internal/usecases/users"
	"go.opentelemetry.io/otel/trace"
)

type UserHandler struct {
	service users.UserService
	tracer  trace.Tracer
}

func NewHandler(service users.UserService, tracer trace.Tracer) (*UserHandler, error) {
	if service == (users.UserService{}) {
		return nil, errors.New("service cannot be empty")
	}
	if tracer == nil {
		return nil, errors.New("tracer cannot be empty")
	}
	return &UserHandler{service, tracer}, nil
}
