package handlers

import (
	"net/http"

	appErrors "github.com/olad5/productive-pulse/pkg/errors"
	response "github.com/olad5/productive-pulse/pkg/utils"
	"github.com/olad5/productive-pulse/todo-service/internal/usecases/todos"
	"github.com/olad5/productive-pulse/todo-service/internal/utils"
)

func (h TodoHandler) GetTodos(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		response.ErrorResponse(w, appErrors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}
	userId, err := h.userService.VerifyUser(ctx, authHeader)
	if err != nil {
		response.ErrorResponse(w, appErrors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	foundTodos, err := h.todoService.GetTodos(ctx, userId)
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
