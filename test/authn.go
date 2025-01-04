package test

import (
	"app/authn"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/fx"
	"golang.org/x/oauth2"
)

var _ authn.UserProvider = (*InMemoryUserProvider)(nil)
var _ authn.TokenProvider = (*JWTTokenProvider)(nil)

var FakeAuthnProviders = fx.Options(
	fx.Provide(NewInMemoryUserProvider),
	fx.Decorate(func(userProvider *InMemoryUserProvider) authn.UserProvider { return userProvider }),
	fx.Decorate(func() authn.TokenProvider { return NewJWTTokenProvider("secret") }),
)

type InMemoryUserProvider struct {
	mu    sync.Mutex
	users map[string]string // email -> password
}

func NewInMemoryUserProvider() *InMemoryUserProvider {
	return &InMemoryUserProvider{
		users: make(map[string]string),
	}
}

func (userProvider *InMemoryUserProvider) Clear() {
	userProvider.mu.Lock()
	defer userProvider.mu.Unlock()

	clear(userProvider.users)
}

func (userProvider *InMemoryUserProvider) Create(ctx context.Context, email string, password string) error {
	userProvider.mu.Lock()
	defer userProvider.mu.Unlock()

	userProvider.users[email] = password

	return nil
}

func (userProvider *InMemoryUserProvider) Delete(ctx context.Context, email string) error {
	userProvider.mu.Lock()
	defer userProvider.mu.Unlock()

	if _, ok := userProvider.users[email]; !ok {
		return authn.ErrUserNotFound
	}

	delete(userProvider.users, email)

	return nil
}

type JWTTokenProvider struct {
	secret string
	exp    time.Duration
}

func NewJWTTokenProvider(secret string) *JWTTokenProvider {
	return &JWTTokenProvider{secret: secret, exp: 7 * 24 * time.Hour}
}

func (tokenProvider *JWTTokenProvider) Get(ctx context.Context, email string, password string) (*oauth2.Token, error) {
	claims := jwt.MapClaims{"sub": email, "email": email}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(tokenProvider.secret))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to sign token", authn.ErrToken)
	}

	return &oauth2.Token{
		TokenType:    "Bearer",
		AccessToken:  tokenString,
		Expiry:       time.Now().Add(tokenProvider.exp),
		RefreshToken: "refresh-token",
	}, nil
}

func (tokenProvider *JWTTokenProvider) Parse(tokenString string) (authn.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: unexpected signing method: %v", authn.ErrToken, token.Header["alg"])
		}

		return []byte(tokenProvider.secret), nil
	})
	if err != nil {
		return authn.Claims{}, authn.ErrParseToken
	}

	if !token.Valid {
		return authn.Claims{}, authn.ErrInvalidSignature
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return authn.Claims{}, authn.ErrParseToken
	}

	return authn.Claims{
		Subject: claims["sub"].(string),
		Email:   claims["email"].(string),
	}, nil
}
