package authn_module

import (
	"app/adapter/casdoor"
	"app/authn"
	"app/authn/authn_http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var (
	AuthnProviderModule = fx.Module(
		"authentication",

		fx.Provide(func(casdoorUserProvider *casdoor.UserProvider) authn.UserProvider {
			return casdoorUserProvider
		}),
		fx.Provide(func(casdoorTokenProvider *casdoor.TokenProvider) authn.TokenProvider {
			return casdoorTokenProvider
		}),
		fx.Provide(authn.NewService),
	)

	HTTPModule = fx.Module(
		"authentication http",

		AuthnProviderModule,

		fx.Provide(authn_http.NewAuthMiddleware),
		fx.Provide(authn_http.NewController, fx.Private),

		fx.Invoke(func(app *fiber.App, authnController *authn_http.Controller) {
			authnController.Register(app)
		}),
	)
)
