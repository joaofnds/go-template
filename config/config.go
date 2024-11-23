package config

import (
	"app/adapter/featureflags"
	"app/adapter/http"
	"app/adapter/logger"
	"app/adapter/metrics"
	"app/adapter/mongo"
	"app/adapter/postgres"
	"app/adapter/redis"
	"app/adapter/tracing"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var (
	Module = fx.Module(
		"config",

		fx.Provide(NewAppConfig),
		fx.Provide(func(config AppConfig) logger.Config { return config.Logger }),
		fx.Provide(func(config AppConfig) http.Config { return config.HTTP }),
		fx.Provide(func(config AppConfig) metrics.Config { return config.Metrics }),
		fx.Provide(func(config AppConfig) postgres.Config { return config.Postgres }),
		fx.Provide(func(config AppConfig) mongo.Config { return config.Mongo }),
		fx.Provide(func(config AppConfig) redis.Config { return config.Redis }),
		fx.Provide(func(config AppConfig) featureflags.Config { return config.FeatureFlags }),
		fx.Provide(func(config AppConfig) tracing.Config { return config.Tracing }),
	)
)

type AppConfig struct {
	Env          string              `mapstructure:"env" validate:"required,oneof=development staging production"`
	Logger       logger.Config       `mapstructure:"logger" validate:"required"`
	HTTP         http.Config         `mapstructure:"http" validate:"required"`
	Metrics      metrics.Config      `mapstructure:"metrics" validate:"required"`
	Postgres     postgres.Config     `mapstructure:"postgres" validate:"required"`
	Mongo        mongo.Config        `mapstructure:"mongo" validate:"required"`
	Redis        redis.Config        `mapstructure:"redis" validate:"required"`
	FeatureFlags featureflags.Config `mapstructure:"feature_flags" validate:"required"`
	Tracing      tracing.Config      `mapstructure:"tracing" validate:"required"`
}

func NewAppConfig(validator *validator.Validate) (AppConfig, error) {
	viperInstance := viper.New()

	viperInstance.MustBindEnv("env", "ENV")
	viperInstance.MustBindEnv("postgres.uri", "POSTGRES_URI")
	viperInstance.MustBindEnv("redis.addr", "REDIS_ADDR")

	viperInstance.SetConfigFile(os.Getenv("CONFIG_PATH"))

	var config AppConfig
	if err := viperInstance.ReadInConfig(); err != nil {
		return config, err
	}

	if err := viperInstance.UnmarshalExact(&config); err != nil {
		return config, err
	}

	return config, validator.Struct(config)
}
