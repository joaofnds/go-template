package redis

type Config struct {
	Addr string `mapstructure:"addr" validate:"required,uri"`
}
