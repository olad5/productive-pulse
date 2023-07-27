package handlers

import (
	"errors"

	"github.com/olad5/productive-pulse/users-service/internal/usecases/users"
)

type Handler struct {
	service users.UserService
}

func NewHandler(service users.UserService) (*Handler, error) {
	if service == (users.UserService{}) {
		return nil, errors.New("service cannot be empty")
	}
	return &Handler{service}, nil
}
