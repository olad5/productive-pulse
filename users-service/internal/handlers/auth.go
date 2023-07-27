package handlers

import (
	"net/http"

	"github.com/olad5/productive-pulse/users-service/internal/utils"
)

func (h Handler) Auth(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")

	userId, err := h.service.VerifyUser(r.Context(), authHeader)
	if err != nil {
		utils.ErrorResponse(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	utils.SuccessResponse(w, "success",
		map[string]interface{}{
			"user_id": userId,
		})
}
