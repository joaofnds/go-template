package tracing

import "time"

type Config struct {
	Addr       string        `mapstructure:"addr" validate:"required,url"`
	Secure     bool          `mapstructure:"secure"`
	Timeout    time.Duration `mapstructure:"timeout" validate:"required"`
	SampleRate float64       `mapstructure:"sample_rate" validate:"required"`
}
