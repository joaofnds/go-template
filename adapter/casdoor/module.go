package casdoor

import (
	"go.uber.org/fx"
)

var (
	Module = fx.Module("authn", Providers, Invokes)

	Providers = fx.Options(
		fx.Provide(NewAuthProvider),
	)
	Invokes = fx.Options()
)
