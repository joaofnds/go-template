package config

import (
	"app/adapter/featureflags"
	"app/adapter/http"
	"app/adapter/metrics"
	"app/adapter/mongo"
	"app/adapter/postgres"
	"app/adapter/redis"
	"os"

	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	Module = fx.Module("config", Providers, Invokes)

	Providers = fx.Options(
		fx.Provide(NewAppConfig),
		fx.Provide(func(config AppConfig) http.Config { return config.HTTP }),
		fx.Provide(func(config AppConfig) metrics.Config { return config.Metrics }),
		fx.Provide(func(config AppConfig) postgres.Config { return config.Postgres }),
		fx.Provide(func(config AppConfig) mongo.Config { return config.Mongo }),
		fx.Provide(func(config AppConfig) redis.Config { return config.Redis }),
		fx.Provide(func(config AppConfig) featureflags.Config { return config.FeatureFlags }),
	)
	Invokes = fx.Options(
		fx.Invoke(LoadConfig),
	)
)

type AppConfig struct {
	Env          string              `mapstructure:"env"`
	HTTP         http.Config         `mapstructure:"http"`
	Metrics      metrics.Config      `mapstructure:"metrics"`
	Postgres     postgres.Config     `mapstructure:"postgres"`
	Mongo        mongo.Config        `mapstructure:"mongo"`
	Redis        redis.Config        `mapstructure:"redis"`
	FeatureFlags featureflags.Config `mapstructure:"feature_flags"`
}

func LoadConfig(logger *zap.Logger) error {
	configFile := os.Getenv("CONFIG_PATH")
	if configFile == "" {
		logger.Info("CONFIG_PATH not set, skipping config file load")
		bindEnvs()
		return nil
	}

	viper.SetConfigFile(configFile)
	return viper.ReadInConfig()
}

func NewAppConfig() (AppConfig, error) {
	var config AppConfig
	return config, viper.UnmarshalExact(&config)
}

func bindEnvs() {
	viper.MustBindEnv("env", "ENV")
	viper.MustBindEnv("metrics.address", "METRICS_ADDRESS")
	viper.MustBindEnv("http.port", "HTTP_PORT")
	viper.MustBindEnv("http.limiter.requests", "HTTP_LIMITER_REQUESTS")
	viper.MustBindEnv("http.limiter.expiration", "HTTP_LIMITER_EXPIRATION")
	viper.MustBindEnv("postgres.uri", "POSTGRES_URI")
	viper.MustBindEnv("mongo.uri", "MONGO_URI")
	viper.MustBindEnv("redis.addr", "REDIS_ADDR")
}
