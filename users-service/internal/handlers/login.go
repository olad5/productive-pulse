package handlers

import (
	"encoding/json"
	"net/http"

	appErrors "github.com/olad5/productive-pulse/pkg/errors"
	response "github.com/olad5/productive-pulse/pkg/utils"
	"github.com/olad5/productive-pulse/users-service/internal/usecases/users"
)

func (u UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := u.tracer.Start(ctx, "login-handler")
	defer span.End()
	if r.Body == nil {
		response.ErrorResponse(w, appErrors.ErrMissingBody, http.StatusBadRequest)
		return
	}
	type requestDTO struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var request requestDTO
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.ErrorResponse(w, appErrors.ErrInvalidJson, http.StatusBadRequest)
		return
	}
	if request.Email == "" {
		response.ErrorResponse(w, "email required", http.StatusBadRequest)
		return
	}
	if request.Password == "" {
		response.ErrorResponse(w, "password required", http.StatusBadRequest)
		return
	}

	accessToken, err := u.service.LogUserIn(ctx, u.tracer, request.Email, request.Password)
	if err != nil && err.Error() == users.ErrUserNotFound {
		response.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil && err.Error() == users.ErrPasswordIncorrect {
		response.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		response.ErrorResponse(w, appErrors.ErrSomethingWentWrong, http.StatusInternalServerError)
		return
	}

	response.SuccessResponse(w, "user logged in successfully",
		map[string]interface{}{
			"access_token": accessToken,
		})
}
