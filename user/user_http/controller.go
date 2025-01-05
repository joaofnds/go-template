package user_http

import (
	"errors"
	"net/http"

	"app/adapter/featureflags"
	"app/authn/authn_http"
	"app/authz/authz_http"
	"app/user"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func NewController(
	validator *validator.Validate,
	authn *authn_http.AuthMiddleware,
	authnz *authz_http.Middleware,
	service *user.Service,
) *Controller {
	return &Controller{
		validator: validator,
		authn:     authn,
		authz:     authnz,
		service:   service,
	}
}

type Controller struct {
	validator *validator.Validate
	authn     *authn_http.AuthMiddleware
	authz     *authz_http.Middleware
	service   *user.Service
}

func (controller *Controller) Register(app *fiber.App) {
	users := app.Group("/users")
	users.Get("/", controller.List)

	user := app.Group(
		"/users/:userID",
		controller.authn.RequireUser,
		controller.middlewareGetUser,
		featureflags.Middleware,
	)
	user.Get(
		"/",
		controller.authz.RequireParamPermission("userID", "user", "read"),
		controller.Get,
	)
	user.Delete(
		"/",
		controller.authz.RequireParamPermission("userID", "user", "delete"),
		controller.Delete,
	)
	user.Get(
		"/feature",
		controller.authz.RequireParamPermission("userID", "user", "get-features"),
		controller.GetFeature,
	)
}

func (controller *Controller) List(ctx *fiber.Ctx) error {
	users, err := controller.service.List(ctx.UserContext())
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.JSON(users)
}

func (controller *Controller) Get(ctx *fiber.Ctx) error {
	return ctx.JSON(ctx.Locals("user"))
}

func (controller *Controller) GetFeature(ctx *fiber.Ctx) error {
	return ctx.JSON(ctx.Locals("flags"))
}

func (controller *Controller) Delete(ctx *fiber.Ctx) error {
	foundUser := ctx.Locals("user").(user.User)

	if err := controller.service.Remove(ctx.UserContext(), foundUser); err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.SendStatus(http.StatusOK)
}

func (controller *Controller) middlewareGetUser(ctx *fiber.Ctx) error {
	userID := ctx.Params("userID")
	ctx.AllParams()
	if userID == "" {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	userFound, err := controller.service.FindByID(ctx.UserContext(), userID)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return ctx.SendStatus(http.StatusNotFound)
		} else {
			return ctx.SendStatus(http.StatusInternalServerError)
		}
	}

	ctx.Locals("user", userFound)
	return ctx.Next()
}
