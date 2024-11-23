package http

import "time"

type Config struct {
	Port    int     `mapstructure:"port" validate:"required"`
	Limiter Limiter `mapstructure:"limiter" validate:"required"`
}

type Limiter struct {
	Requests   int           `mapstructure:"requests" validate:"required"`
	Expiration time.Duration `mapstructure:"expiration" validate:"required"`
}
