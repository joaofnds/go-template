package postgres

type Config struct {
	Addr string `mapstructure:"uri" validate:"required,uri,startswith=postgres://"`
}
