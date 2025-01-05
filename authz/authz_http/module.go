package authz_http

import "go.uber.org/fx"

var Module = fx.Module(
	"authz",

	fx.Provide(NewMiddleware),
)
