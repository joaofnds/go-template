package authn_http

import (
	"app/authn"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	authProvider authn.Provider
}

func NewController(authProvider authn.Provider) *Controller {
	return &Controller{
		authProvider: authProvider,
	}
}

func (controller *Controller) Register(app *fiber.App) {
	auth := app.Group("/auth")
	auth.Post("/login", controller.Login)
}

func (controller *Controller) Login(ctx *fiber.Ctx) error {
	var body LoginBody

	if err := ctx.BodyParser(&body); err != nil {
		return err
	}

	token, err := controller.authProvider.GetToken(body.Username, body.Password)
	if err != nil {
		return err
	}

	return ctx.JSON(token)
}

type LoginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
