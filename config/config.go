package config

import (
	"errors"
	"os"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var Module = fx.Options(fx.Invoke(LoadConfig), fx.Provide(NewAppConfig))

type AppConfig struct {
	Env         string `mapstructure:"env"`
	Port        int    `mapstructure:"port"`
	MongoURI    string `mapstructure:"mongo_uri"`
	MetricsAddr string `mapstructure:"metrics_addr"`
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
