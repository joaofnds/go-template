package test

import (
	"math/rand"
	"web/config"

	"go.uber.org/fx"
)

var RandomAppConfigPort = fx.Decorate(randomAppConfigPort)

func randomAppConfigPort(config config.AppConfig) config.AppConfig {
	config.Port = 10_000 + rand.Intn(5000)
	return config
}
