package handlers

import (
	"net/http"

	appErrors "github.com/olad5/productive-pulse/pkg/errors"
	response "github.com/olad5/productive-pulse/pkg/utils"
)

func (h UserHandler) Auth(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")

	userId, err := h.service.VerifyUser(r.Context(), authHeader)
	if err != nil {
		response.ErrorResponse(w, appErrors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	response.SuccessResponse(w, "success",
		map[string]interface{}{
			"user_id": userId,
		})
}
