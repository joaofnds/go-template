package featureflags

import (
	"context"
	"errors"

	"app/user"

	"github.com/gofiber/fiber/v2"
	ffclient "github.com/thomaspoignant/go-feature-flag"
	"github.com/thomaspoignant/go-feature-flag/ffuser"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"featureflags",
	fx.Invoke(HookFFClient),
)

func HookFFClient(lifecycle fx.Lifecycle, config Config) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return ffclient.Init(ffclient.Config{
				PollingInterval: config.PollingInterval,
				Retriever:       config,
				FileFormat:      "json",
			})
		},
		OnStop: func(ctx context.Context) error {
			ffclient.Close()
			return nil
		},
	})
}

func ForUser(u user.User) (map[string]any, error) {
	flagsUser := ffuser.NewUser(u.Name)

	flags := ffclient.AllFlagsState(flagsUser)
	if !flags.IsValid() {
		return nil, errors.New("invalid flags state")
	}

	allFlags := flags.GetFlags()

	m := make(map[string]any, len(allFlags))
	for name, state := range allFlags {
		m[name] = state.Value
	}
	return m, nil
}

func Middleware(ctx *fiber.Ctx) error {
	if ctx.Locals("user") == nil {
		return ctx.Next()
	}

	flags, err := ForUser(ctx.Locals("user").(user.User))
	if err != nil {
		return fiber.ErrInternalServerError
	}

	ctx.Locals("flags", flags)
	for name, value := range flags {
		ctx.Locals("flags."+name, value)
	}
	return ctx.Next()
}
