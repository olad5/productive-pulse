package handlers

import (
	"errors"

	"github.com/olad5/productive-pulse/todo-service/internal/services/user"
	"github.com/olad5/productive-pulse/todo-service/internal/usecases/todos"
)

type TodoHandler struct {
	todoService todos.TodoService
	userService user.UserServiceAdapter
}

func NewTodoHandler(todoService todos.TodoService, userService user.UserServiceAdapter) (*TodoHandler, error) {
	if todoService == (todos.TodoService{}) {
		return nil, errors.New("TodoService cannot be empty")
	}
	if userService == nil {
		return nil, errors.New("UserService cannot be empty")
	}
	return &TodoHandler{todoService, userService}, nil
}
