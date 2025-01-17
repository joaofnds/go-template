package authn

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/oauth2"
)

var (
	ErrAuthn = errors.New("authn error")

	ErrUser               = fmt.Errorf("user error: %w", ErrAuthn)
	ErrUserNotFound       = fmt.Errorf("user not found: %w", ErrUser)
	ErrFailedToCreateUser = fmt.Errorf("failed to create user: %w", ErrUser)
	ErrFailedToGetUser    = fmt.Errorf("failed to get user: %w", ErrUser)
	ErrFailedToDeleteUser = fmt.Errorf("failed to delete user: %w", ErrUser)

	ErrToken            = fmt.Errorf("token error: %w", ErrAuthn)
	ErrParseToken       = fmt.Errorf("failed to parse token: %w", ErrToken)
	ErrInvalidSignature = fmt.Errorf("signature is invalid: %w", ErrToken)
	ErrParseClaims      = fmt.Errorf("failed to parse claims: %w", ErrToken)
	ErrMissingClaims    = fmt.Errorf("missing claims: %w", ErrToken)
	ErrWrongPassword    = fmt.Errorf("wrong password: %w", ErrToken)
)

type User struct {
	Email string
}

type UserProvider interface {
	Create(ctx context.Context, email string, password string) error
	Delete(ctx context.Context, email string) error
}

type Claims struct {
	Subject string `json:"sub,omitempty"`
	Email   string `json:"email,omitempty"`
}

type TokenProvider interface {
	Get(ctx context.Context, email string, password string) (*oauth2.Token, error)
	Parse(token string) (Claims, error)
}
