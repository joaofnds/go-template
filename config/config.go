package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

const configPath = "CONFIG_PATH"

var (
	Module = fx.Options(fx.Invoke(LoadConfig), fx.Provide(NewAppConfig))
)

type AppConfig struct {
	Env      string `mapstructure:"env"`
	Port     int    `mapstructure:"port"`
	MongoURI string `mapstructure:"mongo_uri"`
}

func LoadConfig() error {
	configFile, ok := os.LookupEnv(configPath)
	if !ok {
		log.Fatalf("could not lookup env %q", configPath)
	}

	viper.SetConfigFile(configFile)
	return viper.ReadInConfig()
}

func NewAppConfig() (AppConfig, error) {
	var config AppConfig

	if err := viper.UnmarshalExact(&config); err != nil {
		return config, err
	}

	return config, nil
}
