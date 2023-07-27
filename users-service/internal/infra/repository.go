package infra

import (
	"context"

	"github.com/olad5/productive-pulse/users-service/internal/domain"
)

var ErrRecordNotFound = "No Record found"

type UserRepository interface {
	Ping(ctx context.Context) error
	CreateUser(ctx context.Context, user domain.User) error
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
}
