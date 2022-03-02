package test

import (
	"math/rand"
	"web/config"
)

func NewAppTestConfig() config.AppConfig {
	return config.AppConfig{
		Env:  "test",
		Port: 10_000 + rand.Intn(5000),
	}
}
