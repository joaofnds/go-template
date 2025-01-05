package authz_http

import (
	"app/authn/authn_http"
	"app/authz"
	"app/internal/ref"
	"app/user"

	"github.com/gofiber/fiber/v2"
)

type Middleware struct {
	enforcer authz.Enforcer
}

func NewMiddleware(enforcer authz.Enforcer) *Middleware {
	return &Middleware{enforcer: enforcer}
}

func (m *Middleware) RequireUserParamPermission(
	param string,
	objectType string,
	action string,
) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		user, ok := ctx.Locals(authn_http.UserKey).(user.User)
		if !ok || user.ID == "" {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}

		objectID := ctx.Params(param)
		if objectID == "" {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		subject := ref.New(user.ID, "user")
		object := ref.New(objectID, objectType)

		if !m.enforcer.Check(authz.NewAppRequest(subject, object, action)) {
			return ctx.SendStatus(fiber.StatusForbidden)
		}

		return ctx.Next()
	}
}
