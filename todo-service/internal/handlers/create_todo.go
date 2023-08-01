package handlers

import (
	"encoding/json"
	"net/http"

	appErrors "github.com/olad5/productive-pulse/pkg/errors"
	response "github.com/olad5/productive-pulse/pkg/utils"
	"github.com/olad5/productive-pulse/todo-service/internal/usecases/todos"
	"github.com/olad5/productive-pulse/todo-service/internal/utils"
)

func (t TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := t.tracer.Start(ctx, "CreateTodo-handler")
	defer span.End()
	if r.Body == nil {
		response.ErrorResponse(w, appErrors.ErrMissingBody, http.StatusBadRequest)
		return
	}

	type requestDTO struct {
		Text string `json:"text"`
	}

	var request requestDTO
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.ErrorResponse(w, appErrors.ErrInvalidJson, http.StatusBadRequest)
		return
	}
	if request.Text == "" {
		response.ErrorResponse(w, "Text required", http.StatusBadRequest)
		return
	}

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
	newTodo, err := t.todoService.CreateTodo(ctx, t.tracer, userId, request.Text)
	if err != nil {
		if err == todos.ErrInvalidUserId {
			response.ErrorResponse(w, appErrors.ErrSomethingWentWrong, http.StatusInternalServerError)
			return
		}
		response.ErrorResponse(w, appErrors.ErrSomethingWentWrong, http.StatusInternalServerError)
		return
	}

	response.SuccessResponse(w, "todo created",
		utils.ToTodoDTO(newTodo))
	return
}
