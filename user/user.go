package user

import (
	"app/internal/ref"
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewRef(id string) ref.Ref {
	return ref.New(id, "user")
}
