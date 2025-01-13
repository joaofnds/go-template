package authn_http

import (
	"app/adapter/validation"
	"app/authn"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	validator *validator.Validate
	auth      *AuthMiddleware
	service   *authn.Service
}

func NewController(
	validator *validator.Validate,
	authMiddleware *AuthMiddleware,
	service *authn.Service,
) *Controller {
	return &Controller{
		validator: validator,
		auth:      authMiddleware,
		service:   service,
	}
}

func (controller *Controller) Register(app *fiber.App) {
	auth := app.Group("/auth")
	auth.
		Get("/userinfo", controller.auth.RequireUser, controller.GetUserInfo).
		Post("/login", controller.Login).
		Post("/register", controller.RegisterUser).
		Delete("/delete", controller.DeleteUser)
}

func (controller *Controller) GetUserInfo(ctx *fiber.Ctx) error {
	return ctx.JSON(ctx.Locals(UserKey))
}

func (controller *Controller) Login(ctx *fiber.Ctx) error {
	var body EmailAndPasswordBody

	if err := ctx.BodyParser(&body); err != nil {
		return err
	}

	token, err := controller.service.Login(ctx.UserContext(), body.Email, body.Password)
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
		return ctx.
			Status(http.StatusBadRequest).
			JSON(fiber.Map{"errors": validation.ErrorMessages(err)})
	}

	createdUser, createUserErr := controller.service.RegisterUser(
		ctx.UserContext(),
		body.Email,
		body.Password,
	)
	if createUserErr != nil {
		return createUserErr
	}

	return ctx.Status(fiber.StatusCreated).JSON(createdUser)
}

func (controller *Controller) DeleteUser(ctx *fiber.Ctx) error {
	email := ctx.Query("email", "")
	if email == "" {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	err := controller.service.DeleteUser(ctx.UserContext(), email)
	if err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

type EmailAndPasswordBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
