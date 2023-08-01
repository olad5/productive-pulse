package todos

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/olad5/productive-pulse/config"
	"github.com/olad5/productive-pulse/todo-service/internal/domain"
	"github.com/olad5/productive-pulse/todo-service/internal/infra"
	"go.opentelemetry.io/otel/trace"
)

type TodoService struct {
	todoRepo infra.TodoRepository

	configurations *config.Configurations
}

var (
	ErrTodoNotFound   = errors.New("record not found")
	ErrInvalidToken   = errors.New("invalid token")
	ErrInvalidTodoId  = errors.New("failing to parse todo uuid")
	ErrInvalidUserId  = errors.New("failing to parse user uuid")
	ErrNotOwnerOfTodo = errors.New("current user is not owner of this todo")
)

func NewTodoService(todoRepo infra.TodoRepository, configurations *config.Configurations) (*TodoService, error) {
	if todoRepo == nil {
		return &TodoService{}, errors.New("TodoService failed to initialize")
	}
	return &TodoService{todoRepo, configurations}, nil
}

func (t *TodoService) CreateTodo(ctx context.Context, tracer trace.Tracer, userId, text string) (domain.Todo, error) {
	ctx, span := tracer.Start(ctx, "CreateTodo-TodoService")
	defer span.End()

	userIdInUUId, err := uuid.Parse(userId)
	if err != nil {
		return domain.Todo{}, ErrInvalidUserId
	}
	newTodo := domain.Todo{
		ID:        uuid.New(),
		UserId:    userIdInUUId,
		Text:      text,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = t.todoRepo.CreateTodo(ctx, newTodo)
	if err != nil {
		return domain.Todo{}, err
	}

	return newTodo, nil
}

func (t *TodoService) UpdateTodo(ctx context.Context, tracer trace.Tracer, userId, existingTodoId, updatedText string) (domain.Todo, error) {
	ctx, span := tracer.Start(ctx, "UpdateTodo-TodoService")
	defer span.End()

	todoIdInUUID, err := uuid.Parse(existingTodoId)
	if err != nil {
		return domain.Todo{}, ErrInvalidTodoId
	}
	userIdInUUID, err := uuid.Parse(userId)
	if err != nil {
		return domain.Todo{}, ErrInvalidUserId
	}

	existingTodo, err := t.todoRepo.GetTodo(ctx, userIdInUUID, todoIdInUUID)
	if err != nil && err.Error() == ErrTodoNotFound.Error() {
		return domain.Todo{}, ErrTodoNotFound
	}

	updatedTodo := domain.Todo{
		ID:        existingTodo.ID,
		UserId:    existingTodo.UserId,
		Text:      updatedText,
		CreatedAt: existingTodo.CreatedAt,
		UpdatedAt: existingTodo.UpdatedAt,
	}

	err = t.todoRepo.UpdateTodo(ctx, updatedTodo)
	if err != nil {
		return domain.Todo{}, err
	}

	return updatedTodo, nil
}

func (t *TodoService) GetTodo(ctx context.Context, tracer trace.Tracer, userId, todoId string) (domain.Todo, error) {
	ctx, span := tracer.Start(ctx, "GetTodo-TodoService")
	defer span.End()

	todoIdInUUID, err := uuid.Parse(todoId)
	if err != nil {
		return domain.Todo{}, ErrInvalidTodoId
	}
	userIdInUUID, err := uuid.Parse(userId)
	if err != nil {
		return domain.Todo{}, ErrInvalidUserId
	}

	todo, err := t.todoRepo.GetTodo(ctx, userIdInUUID, todoIdInUUID)
	if err != nil && err.Error() == ErrTodoNotFound.Error() {
		return domain.Todo{}, ErrTodoNotFound
	}
	if err != nil && err.Error() == ErrNotOwnerOfTodo.Error() {
		return domain.Todo{}, ErrNotOwnerOfTodo
	}
	if err != nil {
		return domain.Todo{}, err
	}

	return todo, nil
}

func (t *TodoService) GetTodos(ctx context.Context, tracer trace.Tracer, userId string) ([]domain.Todo, error) {
	ctx, span := tracer.Start(ctx, "GetTodos-TodoService")
	defer span.End()

	userIdInUUID, err := uuid.Parse(userId)
	if err != nil {
		return []domain.Todo{}, ErrInvalidUserId
	}

	todos, err := t.todoRepo.GetTodos(ctx, userIdInUUID)
	if err != nil {
		return []domain.Todo{}, err
	}

	return todos, nil
}
