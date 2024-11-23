package featureflags

import (
	"context"
	"encoding/json"
	"time"
)

type Config struct {
	PollingInterval time.Duration  `mapstructure:"polling_interval" validate:"required"`
	Flags           map[string]any `mapstructure:"flags" validate:"required,dive,keys,required"`
}

func (c Config) Retrieve(context.Context) ([]byte, error) {
	return json.Marshal(c.Flags)
}
