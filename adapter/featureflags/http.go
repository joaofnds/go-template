package featureflags

import (
	"app/user"

	"github.com/gofiber/fiber/v2"
)

func Middleware(ctx *fiber.Ctx) error {
	u, ok := ctx.Locals("user").(user.User)
	if !ok {
		return ctx.Next()
	}

	flags, err := forKey(u.Name)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	ctx.Locals("flags", flags)
	for name, value := range flags {
		ctx.Locals("flags."+name, value)
	}
	return ctx.Next()
}
