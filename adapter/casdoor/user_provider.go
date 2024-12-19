package casdoor

import (
	"app/authn"
	"context"
	"fmt"

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
		return err
	}

	if !ok {
		return fmt.Errorf("failed to create user")
	}

	return nil
}

func (provider *UserProvider) Delete(ctx context.Context, email string) error {
	casdoorUser, err := provider.casdoor.GetUserByEmail(email)
	if err != nil {
		return err
	}

	if casdoorUser == nil {
		return nil
	}

	ok, err := provider.casdoor.DeleteUser(casdoorUser)
	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("failed to delete user")
	}

	return nil
}
