package validation

import (
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
)

var Module = fx.Provide(validator.New)
