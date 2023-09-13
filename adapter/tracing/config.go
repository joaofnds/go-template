package tracing

import "time"

type Config struct {
	Addr    string        `mapstructure:"addr"`
	Secure  bool          `mapstructure:"secure"`
	Timeout time.Duration `mapstructure:"timeout"`
}
