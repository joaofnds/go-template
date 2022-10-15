package test

import (
	"math/rand"
	"web/http"

	"go.uber.org/fx"
)

var RandomAppConfigPort = fx.Decorate(func(config http.Config) http.Config {
	config.Port = 10_000 + rand.Intn(5000)
	return config
})
