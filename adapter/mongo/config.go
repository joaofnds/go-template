package mongo

type Config struct {
	URI string `mapstructure:"uri" validate:"required,uri,startswith=mongodb"`
}
