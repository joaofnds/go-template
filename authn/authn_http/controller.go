package authn_http

import (
	"app/authn"
	"app/user"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	validator *validator.Validate
	authUsers authn.UserProvider
	tokens    authn.TokenProvider
	users     *user.Service
}

func NewController(
	validator *validator.Validate,
	authUsers authn.UserProvider,
	tokens authn.TokenProvider,
	users *user.Service,
) *Controller {
	return &Controller{
		validator: validator,
		authUsers: authUsers,
		tokens:    tokens,
		users:     users,
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
	var body EmailAndPasswordBody

	if err := ctx.BodyParser(&body); err != nil {
		return err
	}

	token, err := controller.tokens.Get(ctx.UserContext(), body.Email, body.Password)
	if err != nil {
		return err
	}

	return ctx.JSON(token)
}

func (controller *Controller) RegisterUser(ctx *fiber.Ctx) error {
	var body EmailAndPasswordBody

	if err := ctx.BodyParser(&body); err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	if err := controller.validator.Struct(body); err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	createdUser, createUserErr := controller.users.CreateUser(ctx.UserContext(), body.Email)
	if createUserErr != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	createAuthUserErr := controller.authUsers.Create(
		ctx.UserContext(),
		body.Email,
		body.Password,
	)
	if createAuthUserErr != nil {
		return createAuthUserErr
	}

	return ctx.Status(fiber.StatusCreated).JSON(createdUser)
}

func (controller *Controller) DeleteUser(ctx *fiber.Ctx) error {
	emailToDelete := ctx.Query("email", "")
	if emailToDelete == "" {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	err := controller.authUsers.Delete(ctx.UserContext(), emailToDelete)
	if err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

type EmailAndPasswordBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
