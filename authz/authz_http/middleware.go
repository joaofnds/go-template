package authz_http

import (
	"app/authn/authn_http"
	"app/authz"
	"app/internal/ref"
	"app/user"

	"github.com/gofiber/fiber/v2"
)

type Middleware struct {
	permissions authz.PermissionManager
}

func NewMiddleware(permissions authz.PermissionManager) *Middleware {
	return &Middleware{permissions: permissions}
}

func (middleware *Middleware) RequireParamPermission(strRef string, action string) fiber.Handler {
	return middleware.RequirePermission(UserFromLocals, objectFromParams(strRef), action)
}

func (middleware *Middleware) RequirePermission(
	subjectFromContext func(*fiber.Ctx) ref.Ref,
	objectFromContext func(*fiber.Ctx) ref.Ref,
	action string,
) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		subject := subjectFromContext(ctx)
		if subject.ID == "" {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}

		object := objectFromContext(ctx)
		if object.ID == "" {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		if !middleware.permissions.Check(authz.NewRequest(subject, object, action)) {
			return ctx.SendStatus(fiber.StatusForbidden)
		}

		return ctx.Next()
	}
}

func UserFromLocals(ctx *fiber.Ctx) ref.Ref {
	requestUser, ok := ctx.Locals(authn_http.UserKey).(user.User)
	if !ok || requestUser.ID == "" {
		return ref.Ref{}
	}

	return user.NewRef(requestUser.ID)
}

func objectFromParams(strRef string) func(*fiber.Ctx) ref.Ref {
	objRef := ref.NewFromString(strRef)
	return func(ctx *fiber.Ctx) ref.Ref {
		objectID := ctx.Params(objRef.ID)
		if objectID == "" {
			return ref.Ref{}
		}
		return ref.New(objRef.Type, objectID)
	}
}
