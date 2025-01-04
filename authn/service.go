package authn

import (
	"app/user"
	"context"

	"golang.org/x/oauth2"
)

type Service struct {
	tokens TokenProvider
	users  *user.Service
	auth   UserProvider
}

func NewService(
	tokens TokenProvider,
	users *user.Service,
	auth UserProvider,
) *Service {
	return &Service{
		tokens: tokens,
		users:  users,
		auth:   auth,
	}
}

func (service *Service) Login(ctx context.Context, email, password string) (*oauth2.Token, error) {
	return service.tokens.Get(ctx, email, password)
}

func (service *Service) RegisterUser(ctx context.Context, email, password string) (user.User, error) {
	createdUser, createUserErr := service.users.CreateUser(ctx, email)
	if createUserErr != nil {
		return user.User{}, createUserErr
	}

	createAuthUserErr := service.auth.Create(ctx, email, password)
	if createAuthUserErr != nil {
		return user.User{}, createAuthUserErr
	}

	return createdUser, nil
}

func (service *Service) DeleteUser(ctx context.Context, email string) error {
	userToDelete, findErr := service.users.FindByEmail(ctx, email)
	if findErr != nil {
		return findErr
	}

	deleteErr := service.users.Remove(ctx, userToDelete)
	if deleteErr != nil {
		return deleteErr
	}

	return service.auth.Delete(ctx, email)
}
