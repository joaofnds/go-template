package health

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewHealthService),
	fx.Provide(func(service *Service) Checker { return service }),
)

const (
	StatusUp   = "up"
	StatusDown = "down"
)

type Check struct {
	Mongo Status `json:"mongo"`
	Redis Status `json:"redis"`
}

func (c Check) AllUp() bool {
	return c.Mongo.IsUp() && c.Redis.IsUp()
}

type Status struct {
	Status string `json:"status"`
}

func (s Status) IsUp() bool {
	return s.Status == StatusUp
}
