package handlers

import (
	"errors"

	"github.com/olad5/productive-pulse/users-service/internal/usecases/users"
)

type UserHandler struct {
	service users.UserService
}

func NewHandler(service users.UserService) (*UserHandler, error) {
	if service == (users.UserService{}) {
		return nil, errors.New("service cannot be empty")
	}
	return &UserHandler{service}, nil
}
