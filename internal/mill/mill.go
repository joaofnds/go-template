package mill

import (
	"app/internal/util"
	"context"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

func NewEventHandler[T any](handlerFunc func(ctx context.Context, event *T) error) cqrs.EventHandler {
	return cqrs.NewEventHandler(util.FunctionName(handlerFunc), handlerFunc)
}
