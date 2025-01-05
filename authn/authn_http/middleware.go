package authn_http

import (
	"app/authn"
	"app/user"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

const (
	UserKey = "user"
)

type AuthMiddleware struct {
	tokens authn.TokenProvider
	users  *user.Service
}

func NewAuthMiddleware(
	token authn.TokenProvider,
	users *user.Service,
) *AuthMiddleware {
	return &AuthMiddleware{
		tokens: token,
		users:  users,
	}
}

func (authMiddleware *AuthMiddleware) RequireUser(ctx *fiber.Ctx) error {
	authorization := ctx.Get("Authorization")
	if authorization == "" {
		return ctx.SendStatus(http.StatusUnauthorized)
	}

	claims, err := authMiddleware.tokens.Parse(strings.TrimPrefix(authorization, "Bearer "))
	if err != nil {
		return ctx.SendStatus(http.StatusUnauthorized)
	}

	user, err := authMiddleware.users.FindByEmail(ctx.UserContext(), claims.Email)
	if err != nil {
		return ctx.SendStatus(http.StatusUnauthorized)
	}

	ctx.Locals(UserKey, user)

	return ctx.Next()
}
