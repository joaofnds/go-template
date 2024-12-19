package authn_http

import (
	"app/adapter/casdoor"
	"app/authn"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"authentication",

	fx.Provide(func(casdoorUserProvider *casdoor.UserProvider) authn.UserProvider {
		return casdoorUserProvider
	}, fx.Private),
	fx.Provide(func(casdoorTokenProvider *casdoor.TokenProvider) authn.TokenProvider {
		return casdoorTokenProvider
	}, fx.Private),
	fx.Provide(NewController),
)
