package http

import "time"

type Config struct {
	Port    int     `mapstructure:"port"`
	Limiter Limiter `mapstructure:"limiter"`
}

type Limiter struct {
	Requests   int           `mapstructure:"requests"`
	Expiration time.Duration `mapstructure:"expiration"`
}
