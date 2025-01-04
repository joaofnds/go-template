package authn

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/oauth2"
)

var (
	ErrAuthn = errors.New("authn error")

	ErrUser               = fmt.Errorf("%w: user error", ErrAuthn)
	ErrUserNotFound       = fmt.Errorf("%w: user not found", ErrUser)
	ErrFailedToCreateUser = fmt.Errorf("%w: failed to create user", ErrUser)
	ErrFailedToGetUser    = fmt.Errorf("%w: failed to get user", ErrUser)
	ErrFailedToDeleteUser = fmt.Errorf("%w: failed to delete user", ErrUser)

	ErrToken            = fmt.Errorf("%w: token error", ErrAuthn)
	ErrParseToken       = fmt.Errorf("%w: failed to parse token", ErrToken)
	ErrInvalidSignature = fmt.Errorf("%w: signature is invalid", ErrToken)
	ErrParseClaims      = fmt.Errorf("%w: failed to parse claims", ErrToken)
	ErrMissingClaims    = fmt.Errorf("%w: missing claims", ErrToken)
	ErrWrongPassword    = fmt.Errorf("%w: wrong password", ErrToken)
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
