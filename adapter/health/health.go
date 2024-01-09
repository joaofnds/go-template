package health

import "go.uber.org/fx"

var Module = fx.Module(
	"health",
	fx.Provide(NewHealthService),
	fx.Provide(func(service *Service) Checker { return service }),
	fx.Provide(NewHealthController),
)
