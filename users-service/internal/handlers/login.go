package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/olad5/productive-pulse/users-service/internal/usecases/users"
	"github.com/olad5/productive-pulse/users-service/internal/utils"
)

func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		utils.ErrorResponse(w, "missing body request", http.StatusBadRequest)
		return
	}
	type requestDTO struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var request requestDTO
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		utils.ErrorResponse(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if request.Email == "" {
		utils.ErrorResponse(w, "email required", http.StatusBadRequest)
		return
	}
	if request.Password == "" {
		utils.ErrorResponse(w, "password required", http.StatusBadRequest)
		return
	}

	accessToken, err := h.service.LogUserIn(r.Context(), request.Email, request.Password)
	if err != nil && err.Error() == users.ErrUserNotFound {
		utils.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil && err.Error() == users.ErrPasswordIncorrect {
		utils.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		utils.ErrorResponse(w, "something went wrong", http.StatusBadRequest)
		return
	}

	utils.SuccessResponse(w, "user logged in successfully",
		map[string]interface{}{
			"access_token": accessToken,
		})
}
