package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/trace"
)

type UserServiceResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		UserId string `json:"user_id"`
	} `json:"data"`
}

type UserServiceAdapter interface {
	VerifyUser(ctx context.Context, tracer trace.Tracer, authHeader string) (string, error)
}

type UserService struct {
	client *http.Client
	url    string
}

type UserId string

func NewUserService(client *http.Client, url string) (*UserService, error) {
	if client == nil {
		return nil, errors.New("client cannot be nil")
	}
	if url == "" {
		return nil, errors.New("url cannot be empty")
	}
	return &UserService{client, url}, nil
}

func (u *UserService) VerifyUser(ctx context.Context, tracer trace.Tracer, authHeader string) (string, error) {
	ctx, span := tracer.Start(ctx, "user-service-adapter")
	defer span.End()

	url := u.url + "/users/auth"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", authHeader)

	res, err := u.client.Do(req.WithContext(ctx))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad request to UserService: %d", res.StatusCode)
	}

	var ur UserServiceResponse
	if err := json.NewDecoder(res.Body).Decode(&ur); err != nil {
		return "", fmt.Errorf("could not decode the response body of UserService: %w", err)
	}

	return ur.Data.UserId, nil
}
