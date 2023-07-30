package domain

import (
	"time"

	"github.com/google/uuid"
)

type Todo struct {
	ID        uuid.UUID
	UserId    uuid.UUID
	Text      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
