package http

import (
	"web/http/fiber"
	"web/http/health"
	"web/http/user"

	"go.uber.org/fx"
)

var (
	Module = fx.Options(
		fiber.Module,
		user.Providers,
		health.Providers,
	)
)
