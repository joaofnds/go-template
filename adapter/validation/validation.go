package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
)

var Module = fx.Provide(validator.New)

func ErrorMessages(err error) []string {
	errors, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil
	}

	var errorMessages []string
	for _, err := range errors {
		errorMessages = append(
			errorMessages,
			fmt.Sprintf(`Field validation for '%s' failed on the '%s' tag`, err.Field(), err.Tag()),
		)
	}
	return errorMessages
}
