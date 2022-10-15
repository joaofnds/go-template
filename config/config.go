package config

import (
	"errors"
	"os"
	"web/mongo"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Invoke(LoadConfig),
	fx.Provide(NewAppConfig),
	fx.Provide(func(config AppConfig) mongo.Config { return config.Mongo }),
)

type AppConfig struct {
	Env         string       `mapstructure:"env"`
	Port        int          `mapstructure:"port"`
	Mongo       mongo.Config `mapstructure:"mongo"`
	MetricsAddr string       `mapstructure:"metrics_addr"`
}

func LoadConfig() error {
	configFile := os.Getenv("CONFIG_PATH")
	if configFile == "" {
		return errors.New("CONFIG_PATH env not set")
	}

	viper.SetConfigFile(configFile)
	return viper.ReadInConfig()
}

func NewAppConfig() (AppConfig, error) {
	var config AppConfig
	return config, viper.UnmarshalExact(&config)
}
