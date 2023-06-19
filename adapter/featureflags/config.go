package featureflags

import (
	"context"
	"encoding/json"
	"time"
)

type Config struct {
	PollingInterval time.Duration  `mapstructure:"polling_interval"`
	Flags           map[string]any `mapstructure:"flags"`
}

func (c Config) Retrieve(context.Context) ([]byte, error) {
	return json.Marshal(c.Flags)
}
