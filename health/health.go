package health

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewHealthService),
	fx.Provide(func(healthService *HealthService) HealthChecker {
		return healthService
	}),
)

const (
	StatusUp   = "up"
	StatusDown = "down"
)

type HealthCheck struct {
	Mongo Status `json:"mongo"`
	Redis Status `json:"redis"`
}

func (hc HealthCheck) AllUp() bool {
	return hc.Mongo.IsUp() && hc.Redis.IsUp()
}

type Status struct {
	Status string `json:"status"`
}

func (s Status) IsUp() bool {
	return s.Status == StatusUp
}
