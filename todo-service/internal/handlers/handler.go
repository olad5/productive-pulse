package handlers

import (
	"errors"

	"github.com/olad5/productive-pulse/todo-service/internal/services/user"
	"github.com/olad5/productive-pulse/todo-service/internal/usecases/todos"
	"go.opentelemetry.io/otel/trace"
)

type TodoHandler struct {
	todoService todos.TodoService
	userService user.UserServiceAdapter
	tracer      trace.Tracer
}

func NewTodoHandler(todoService todos.TodoService, userService user.UserServiceAdapter, tracer trace.Tracer) (*TodoHandler, error) {
	if todoService == (todos.TodoService{}) {
		return nil, errors.New("TodoService cannot be empty")
	}
	if userService == nil {
		return nil, errors.New("UserService cannot be empty")
	}
	if tracer == nil {
		return nil, errors.New("tracer cannot be empty")
	}
	return &TodoHandler{todoService, userService, tracer}, nil
}
