package integration

import (
	"context"
	"errors"
	"net/http"
)

type StubUserService struct {
	client *http.Client
	url    string
}

func (s *StubUserService) VerifyUser(ctx context.Context, authHeader string) (string, error) {
	if authHeader == "Bearer "+ValidTokenForUser1 {
		return "8aa74031-393d-401e-a5d3-ba72089abe40", nil
	}
	if authHeader == "Bearer "+ValidTokenForUser2 {
		return "9a98dc85-fe0a-4cbb-8bb7-f67fceae7751", nil
	}
	return "", errors.New("error decoding jwt")
}
