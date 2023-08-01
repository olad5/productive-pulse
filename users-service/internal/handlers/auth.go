package handlers

import (
	"net/http"

	appErrors "github.com/olad5/productive-pulse/pkg/errors"
	response "github.com/olad5/productive-pulse/pkg/utils"
)

func (u UserHandler) Auth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := u.tracer.Start(ctx, "auth-handler")
	defer span.End()

	authHeader := r.Header.Get("Authorization")

	userId, err := u.service.VerifyUser(ctx, u.tracer, authHeader)
	if err != nil {
		response.ErrorResponse(w, appErrors.ErrUnauthorized, http.StatusUnauthorized)
		return
	}

	response.SuccessResponse(w, "success",
		map[string]interface{}{
			"user_id": userId,
		})
}
