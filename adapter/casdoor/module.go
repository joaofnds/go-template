package casdoor

import (
	"crypto/rsa"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/fx"
)

var (
	Module = fx.Module("authn", Providers, Invokes)

	Providers = fx.Options(
		fx.Provide(NewCertificate, fx.Private),
		fx.Provide(NewClient, fx.Private),
		fx.Provide(NewTokenProvider),
		fx.Provide(NewUserProvider),
	)
	Invokes = fx.Options()
)

func NewCertificate(config Config) (*rsa.PublicKey, error) {
	return jwt.ParseRSAPublicKeyFromPEM([]byte(config.Certificate))
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
