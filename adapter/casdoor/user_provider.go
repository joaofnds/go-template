package casdoor

import (
	"app/authn"
	"context"
	"fmt"
	"strings"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
)

var _ authn.UserProvider = &UserProvider{}

type UserProvider struct {
	config  Config
	casdoor *casdoorsdk.Client
}

func NewUserProvider(
	config Config,
	casdoor *casdoorsdk.Client,
) *UserProvider {
	return &UserProvider{
		config:  config,
		casdoor: casdoor,
	}
}

func (provider *UserProvider) Create(
	ctx context.Context,
	email string,
	password string,
) error {
	ok, err := provider.casdoor.AddUser(&casdoorsdk.User{
		Owner:    provider.config.OrganizationName,
		Name:     email,
		Email:    email,
		Password: password,
	})
	if err != nil {
		return fmt.Errorf("%w: %s", authn.ErrFailedToCreateUser, err)
	}

	if !ok {
		return fmt.Errorf("%w: %s", authn.ErrFailedToCreateUser, email)
	}

	return nil
}

func (provider *UserProvider) Delete(ctx context.Context, email string) error {
	casdoorUser, getErr := provider.getUser(email)
	if getErr != nil {
		return getErr
	}

	if casdoorUser == nil {
		return nil
	}

	if _, deleteErr := provider.casdoor.DeleteUser(casdoorUser); deleteErr != nil {
		return fmt.Errorf("%w: %s", authn.ErrFailedToDeleteUser, deleteErr)
	}

	return nil
}

func (provider *UserProvider) getUser(email string) (*casdoorsdk.User, error) {
	casdoorUser, err := provider.casdoor.GetUserByEmail(email)

	switch {
	case err == nil:
		return casdoorUser, nil
	case strings.Contains(err.Error(), "doesn't exist"):
		return casdoorUser, fmt.Errorf("%w: %s", authn.ErrUserNotFound, email)
	default:
		return casdoorUser, fmt.Errorf("%w: %s", authn.ErrFailedToGetUser, err)
	}
}
