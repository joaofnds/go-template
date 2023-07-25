package featureflags

import (
	"context"
	"errors"

	ffclient "github.com/thomaspoignant/go-feature-flag"
	"github.com/thomaspoignant/go-feature-flag/ffcontext"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"featureflags",
	fx.Invoke(HookClient),
)

func HookClient(lifecycle fx.Lifecycle, config Config) {
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

func forKey(key string) (map[string]any, error) {
	flags := ffclient.AllFlagsState(ffcontext.NewEvaluationContext(key))
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
