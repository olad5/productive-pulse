package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/olad5/productive-pulse/users-service/internal/usecases/users"
	"github.com/olad5/productive-pulse/users-service/internal/utils"
)

func (h Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		utils.ErrorResponse(w, "missing body request", http.StatusBadRequest)
		return
	}
	type requestDTO struct {
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Password  string `json:"password"`
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

	if request.FirstName == "" {
		utils.ErrorResponse(w, "first_name required", http.StatusBadRequest)
		return
	}
	if request.LastName == "" {
		utils.ErrorResponse(w, "last_name required", http.StatusBadRequest)
		return
	}

	newUser, err := h.service.CreateUser(r.Context(), request.FirstName, request.LastName, request.Email, request.Password)
	if err != nil && err.Error() == users.ErrUserAlreadyExists {
		utils.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		utils.ErrorResponse(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, "user created successfully",
		map[string]interface{}{
			"id":         newUser.ID.String(),
			"email":      newUser.Email,
			"first_name": newUser.FirstName,
			"last_name":  newUser.LastName,
		})
}
