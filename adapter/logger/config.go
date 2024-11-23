package logger

type Config struct {
	Level string `json:"level" validate:"required"`
}
