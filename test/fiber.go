package test

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var TestFiber = fx.Provide(func() *fiber.App { return fiber.New() })
