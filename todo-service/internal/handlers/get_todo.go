package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	appErrors "github.com/olad5/productive-pulse/pkg/errors"
	response "github.com/olad5/productive-pulse/pkg/utils"
	"github.com/olad5/productive-pulse/todo-service/internal/usecases/todos"
	"github.com/olad5/productive-pulse/todo-service/internal/utils"
)

func (h TodoHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	todoId := chi.URLParam(r, "id")

	if todoId == "" {
		response.ErrorResponse(w, "todo id required", http.StatusBadRequest)
		return
	}

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

	todo, err := h.todoService.GetTodo(ctx, userId, todoId)
	if err != nil {
		if err == todos.ErrInvalidTodoId {
			response.ErrorResponse(w, "invalid todoId", http.StatusBadRequest)
			return
		}
		if err == todos.ErrTodoNotFound {
			response.ErrorResponse(w, err.Error(), http.StatusNotFound)
			return
		}

		if err == todos.ErrNotOwnerOfTodo {
			response.ErrorResponse(w, appErrors.ErrUnauthorized, http.StatusUnauthorized)
			return
		}
		response.ErrorResponse(w, appErrors.ErrSomethingWentWrong, http.StatusInternalServerError)
		return
	}

	response.SuccessResponse(w, "todo retreived",
		utils.ToTodoDTO(todo))
	return
}
