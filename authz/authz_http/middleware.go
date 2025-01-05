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

func (middleware *Middleware) RequireParamPermission(
	key string,
	objectType string,
	action string,
) fiber.Handler {
	return middleware.RequirePermission(
		UserFromLocals,
		objectFromParams(key, objectType),
		action,
	)
}

func (middleware *Middleware) RequirePermission(
	subjectFromContext func(*fiber.Ctx) ref.Ref,
	objectRef func(*fiber.Ctx) ref.Ref,
	action string,
) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		subject := subjectFromContext(ctx)
		if subject.ID == "" {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}

		object := objectRef(ctx)
		if object.ID == "" {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		if !middleware.enforcer.Check(authz.NewAppRequest(subject, object, action)) {
			return ctx.SendStatus(fiber.StatusForbidden)
		}

		return ctx.Next()
	}
}

func UserFromLocals(ctx *fiber.Ctx) ref.Ref {
	user, ok := ctx.Locals(authn_http.UserKey).(user.User)
	if !ok || user.ID == "" {
		return ref.Ref{}
	}

	return ref.NewUser(user.ID)
}

func objectFromParams(key, objectType string) func(*fiber.Ctx) ref.Ref {
	return func(ctx *fiber.Ctx) ref.Ref {
		objectID := ctx.Params(key)
		if objectID == "" {
			return ref.Ref{}
		}

		return ref.New(objectID, objectType)
	}
}
