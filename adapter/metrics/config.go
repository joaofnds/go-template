package metrics

type Config struct {
	Addr string `mapstructure:"addr" validate:"required,hostname_port"`
}
