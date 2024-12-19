package authn

import (
	"context"

	"golang.org/x/oauth2"
)

type UserProvider interface {
	Create(ctx context.Context, email string, password string) error
	Delete(ctx context.Context, email string) error
}

type TokenProvider interface {
	Get(ctx context.Context, email string, password string) (*oauth2.Token, error)
	Parse(token string) (Claims, error)
}

type Claims struct {
	Subject string `json:"sub,omitempty"`
	Email   string `json:"email,omitempty"`
}
