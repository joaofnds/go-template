package casdoor

import (
	"app/authn"
	"context"
	"crypto/rsa"
	"fmt"

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
	return provider.oauth.PasswordCredentialsToken(ctx, email, password)
}

func (provider *TokenProvider) Parse(rawToken string) (authn.Claims, error) {
	claims := authn.Claims{}

	parsedToken, parseErr := jwt.Parse(
		rawToken,
		func(_ *jwt.Token) (any, error) { return provider.publicKey, nil },
		jwt.WithValidMethods([]string{"RS256"}),
	)

	if parseErr != nil {
		return claims, parseErr
	}

	if !parsedToken.Valid {
		return claims, jwt.ErrSignatureInvalid
	}

	tokenClaims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return claims, fmt.Errorf("failed to parse token claims")
	}

	if sub, ok := tokenClaims["sub"].(string); !ok {
		return claims, fmt.Errorf("missing subject claim")
	} else {
		claims.Subject = sub
	}

	if email, ok := tokenClaims["email"].(string); !ok {
		return claims, fmt.Errorf("missing email claim")
	} else {
		claims.Email = email
	}

	return claims, nil
}
