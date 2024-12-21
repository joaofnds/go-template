package casdoor

import (
	"crypto/rsa"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/fx"
	"golang.org/x/oauth2"
)

var (
	Module = fx.Module("authn", Providers, Invokes)

	Providers = fx.Options(
		fx.Provide(NewCertificate, fx.Private),
		fx.Provide(NewOAuth2Config, fx.Private),
		fx.Provide(NewClient, fx.Private),
		fx.Provide(NewTokenProvider),
		fx.Provide(NewUserProvider),
	)
	Invokes = fx.Options()
)

func NewCertificate(config Config) (*rsa.PublicKey, error) {
	return jwt.ParseRSAPublicKeyFromPEM([]byte(config.Certificate))
}

func NewOAuth2Config(config Config) oauth2.Config {
	return oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:   config.Endpoint + "/api/login/oauth/authorize",
			TokenURL:  config.Endpoint + "/api/login/oauth/access_token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
		Scopes: nil,
	}
}

func NewClient(config Config) *casdoorsdk.Client {
	return casdoorsdk.NewClient(
		config.Endpoint,
		config.ClientID,
		config.ClientSecret,
		config.Certificate,
		config.OrganizationName,
		config.ApplicationName,
	)
}
