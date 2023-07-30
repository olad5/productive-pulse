package infra

import (
	"context"

	"github.com/google/uuid"
	"github.com/olad5/productive-pulse/todo-service/internal/domain"
)

type TodoRepository interface {
	Ping(ctx context.Context) error
	CreateTodo(ctx context.Context, todo domain.Todo) error
	UpdateTodo(ctx context.Context, todo domain.Todo) error
	GetTodo(ctx context.Context, userId, todoId uuid.UUID) (domain.Todo, error)
	GetTodos(ctx context.Context, userId uuid.UUID) ([]domain.Todo, error)
}
