package test

import (
	"math/rand"
	"web/config"
)

func RandomAppConfigPort(config config.AppConfig) config.AppConfig {
	config.Port = 10_000 + rand.Intn(5000)
	return config
}
