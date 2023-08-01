package handlers

import (
	"net/http"

	appErrors "github.com/olad5/productive-pulse/pkg/errors"
	response "github.com/olad5/productive-pulse/pkg/utils"
	"github.com/olad5/productive-pulse/todo-service/internal/usecases/todos"
	"github.com/olad5/productive-pulse/todo-service/internal/utils"
)

func (t TodoHandler) GetTodos(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := t.tracer.Start(ctx, "GetTodos-handler")
	defer span.End()

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		response.ErrorResponse(w, appErrors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}
	userId, err := t.userService.VerifyUser(ctx, t.tracer, authHeader)
	if err != nil {
		response.ErrorResponse(w, appErrors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	foundTodos, err := t.todoService.GetTodos(ctx, t.tracer, userId)
	if err != nil && err == todos.ErrInvalidUserId {
		response.ErrorResponse(w, appErrors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}
	if err != nil {
		response.ErrorResponse(w, appErrors.ErrSomethingWentWrong, http.StatusInternalServerError)
		return
	}

	var todosData []map[string]interface{}
	for _, todo := range foundTodos {
		todoData := utils.ToTodoDTO(todo)
		todosData = append(todosData, todoData)
	}

	response.SuccessResponse(w, "todos retreived", todosData)
}
