package mill

import (
	"app/internal/util"
	"context"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

func NewEventHandler[T any](handlerFunc func(ctx context.Context, event *T) error) cqrs.EventHandler {
	return cqrs.NewEventHandler(util.FunctionName(handlerFunc), handlerFunc)
}

type SelfRegisterHandler interface {
	RegisterEventHandlers(processor *cqrs.EventProcessor) error
}

func RegisterEventHandlers(processor *cqrs.EventProcessor, handlers ...SelfRegisterHandler) error {
	for _, handler := range handlers {
		if err := handler.RegisterEventHandlers(processor); err != nil {
			return err
		}
	}
	return nil
}
