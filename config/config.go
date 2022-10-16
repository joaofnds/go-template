package config

import (
	"app/http"
	"app/metrics"
	"app/mongo"
	"os"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Invoke(LoadConfig),
	fx.Provide(NewAppConfig),
	fx.Provide(func(config AppConfig) http.Config { return config.HTTP }),
	fx.Provide(func(config AppConfig) mongo.Config { return config.Mongo }),
	fx.Provide(func(config AppConfig) metrics.Config { return config.Metrics }),
)

type AppConfig struct {
	Env     string         `mapstructure:"env"`
	HTTP    http.Config    `mapstructure:"http"`
	Mongo   mongo.Config   `mapstructure:"mongo"`
	Metrics metrics.Config `mapstructure:"metrics"`
}

func init() {
	viper.MustBindEnv("env", "ENV")
	viper.MustBindEnv("metrics.address", "METRICS_ADDRESS")
	viper.MustBindEnv("http.port", "HTTP_PORT")
	viper.MustBindEnv("http.limiter.requests", "HTTP_LIMITER_REQUESTS")
	viper.MustBindEnv("http.limiter.expiration", "HTTP_LIMITER_EXPIRATION")
	viper.MustBindEnv("mongo.uri", "MONGO_URI")
}

func LoadConfig() error {
	configFile := os.Getenv("CONFIG_PATH")
	if configFile == "" {
		return nil
	}

	viper.SetConfigFile(configFile)
	return viper.ReadInConfig()
}

func NewAppConfig() (AppConfig, error) {
	var config AppConfig
	return config, viper.UnmarshalExact(&config)
}
