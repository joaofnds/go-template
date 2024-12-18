package casdoor

import (
	"app/authn"
	"context"
	"crypto/rsa"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

var _ authn.Provider = &AuthProvider{}

type AuthProvider struct {
	config    Config
	publicKey *rsa.PublicKey
}

func NewAuthProvider(config Config) (*AuthProvider, error) {
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(config.Certificate))
	if err != nil {
		return nil, err
	}

	return &AuthProvider{
		config:    config,
		publicKey: publicKey,
	}, nil
}

func (provider *AuthProvider) GetToken(username string, password string) (*oauth2.Token, error) {
	config := oauth2.Config{
		ClientID:     provider.config.ClientID,
		ClientSecret: provider.config.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:   fmt.Sprintf("%s/api/login/oauth/authorize", provider.config.Endpoint),
			TokenURL:  fmt.Sprintf("%s/api/login/oauth/access_token", provider.config.Endpoint),
			AuthStyle: oauth2.AuthStyleInParams,
		},
		Scopes: nil,
	}

	return config.PasswordCredentialsToken(context.Background(), username, password)
}

func (provider *AuthProvider) ParseToken(rawToken string) (authn.Claims, error) {
	claims := authn.Claims{}

	parsedToken, parseErr := jwt.Parse(
		rawToken,
		func(_ *jwt.Token) (interface{}, error) { return provider.publicKey, nil },
		jwt.WithValidMethods([]string{"RS256"}),
	)
	if parseErr != nil {
		return claims, parseErr
	}

	if subject, err := parsedToken.Claims.GetSubject(); err != nil {
		return claims, err
	} else {
		claims.Subject = subject
	}

	return claims, nil
}
