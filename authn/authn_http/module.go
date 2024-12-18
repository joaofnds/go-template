package authn_http

import (
	"app/adapter/casdoor"
	"app/authn"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"authentication",

	fx.Provide(func(casdoorAuthProvider *casdoor.AuthProvider) authn.Provider {
		return casdoorAuthProvider
	}, fx.Private),
	fx.Provide(NewController),
)
