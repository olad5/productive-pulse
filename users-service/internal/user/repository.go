package user

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	Ping(ctx context.Context) error
}

type PostgresRepository struct {
	connection *pgx.Conn
}

func NewPostgresRepo(ctx context.Context, DatabaseUrl string) (*PostgresRepository, error) {
	conn, err := pgx.Connect(ctx, DatabaseUrl)
	if err != nil {
		return &PostgresRepository{}, fmt.Errorf("Failed to create PostgresRepository:  %w", err)
	}
	return &PostgresRepository{connection: conn}, nil
}

func (p *PostgresRepository) Ping(ctx context.Context) error {
	err := p.connection.Ping(ctx)
	if err != nil {
		return fmt.Errorf("Failed to Ping Postgres DB:  %w", err)
	}

	return nil
}
