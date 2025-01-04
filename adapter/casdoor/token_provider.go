package casdoor

import (
	"app/authn"
	"context"
	"crypto/rsa"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

var _ authn.TokenProvider = &TokenProvider{}

type TokenProvider struct {
	config    Config
	publicKey *rsa.PublicKey
	oauth     oauth2.Config
}

func NewTokenProvider(
	config Config,
	oauth oauth2.Config,
	publicKey *rsa.PublicKey,
) *TokenProvider {
	return &TokenProvider{
		config:    config,
		oauth:     oauth,
		publicKey: publicKey,
	}
}

func (provider *TokenProvider) Get(
	ctx context.Context,
	email string,
	password string,
) (*oauth2.Token, error) {
	token, err := provider.oauth.PasswordCredentialsToken(ctx, email, password)
	switch {
	case err == nil:
		return token, nil
	case strings.Contains(err.Error(), "the user does not exist"):
		return nil, authn.ErrUserNotFound
	case strings.Contains(err.Error(), "invalid username or password"):
		return nil, authn.ErrWrongPassword
	default:
		return nil, fmt.Errorf("%w: %s", authn.ErrFailedToGetUser, err)
	}
}

func (provider *TokenProvider) Parse(rawToken string) (authn.Claims, error) {
	claims := authn.Claims{}

	parsedToken, parseErr := jwt.Parse(
		rawToken,
		func(_ *jwt.Token) (any, error) { return provider.publicKey, nil },
		jwt.WithValidMethods([]string{"RS256"}),
	)

	if parseErr != nil {
		return claims, fmt.Errorf("%w: %s", authn.ErrParseToken, parseErr)
	}

	if !parsedToken.Valid {
		return claims, authn.ErrInvalidSignature
	}

	tokenClaims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return claims, authn.ErrParseClaims
	}

	if sub, ok := tokenClaims["sub"].(string); !ok {
		return claims, fmt.Errorf("%w: missing sub claim", authn.ErrMissingClaims)
	} else {
		claims.Subject = sub
	}

	if email, ok := tokenClaims["email"].(string); !ok {
		return claims, fmt.Errorf("%w: missing email claim", authn.ErrMissingClaims)
	} else {
		claims.Email = email
	}

	return claims, nil
}
