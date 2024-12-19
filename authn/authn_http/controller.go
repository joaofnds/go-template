package authn_http

import (
	"app/authn"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	users  authn.UserProvider
	tokens authn.TokenProvider
}

func NewController(
	users authn.UserProvider,
	tokens authn.TokenProvider,
) *Controller {
	return &Controller{
		users:  users,
		tokens: tokens,
	}
}

func (controller *Controller) Register(app *fiber.App) {
	auth := app.Group("/auth")
	auth.
		Get("/userinfo", controller.GetUserInfo).
		Post("/login", controller.Login).
		Post("/register", controller.RegisterUser).
		Delete("/delete", controller.DeleteUser)
}

func (controller *Controller) GetUserInfo(ctx *fiber.Ctx) error {
	authorization := ctx.Get("Authorization")
	if authorization == "" {
		return ctx.SendStatus(http.StatusUnauthorized)
	}

	claims, err := controller.tokens.Parse(strings.TrimPrefix(authorization, "Bearer "))
	if err != nil {
		return ctx.SendStatus(http.StatusUnauthorized)
	}

	return ctx.JSON(claims)
}

func (controller *Controller) Login(ctx *fiber.Ctx) error {
	var body UsernamePasswordBody

	if err := ctx.BodyParser(&body); err != nil {
		return err
	}

	token, err := controller.tokens.Get(ctx.UserContext(), body.Username, body.Password)
	if err != nil {
		return err
	}

	return ctx.JSON(token)
}

func (controller *Controller) RegisterUser(ctx *fiber.Ctx) error {
	var body UsernamePasswordBody

	if err := ctx.BodyParser(&body); err != nil {
		return err
	}

	err := controller.users.Create(
		ctx.UserContext(),
		body.Username,
		body.Password,
	)
	if err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusCreated)
}

func (controller *Controller) DeleteUser(ctx *fiber.Ctx) error {
	emailToDelete := ctx.Query("email", "")
	if emailToDelete == "" {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	err := controller.users.Delete(ctx.UserContext(), emailToDelete)
	if err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

type UsernamePasswordBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
