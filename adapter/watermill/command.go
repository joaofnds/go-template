package watermill

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
)

func newCommandBus(
	logger watermill.LoggerAdapter,
	marshaler cqrs.CommandEventMarshaler,
	publisher message.Publisher,
) (*cqrs.CommandBus, error) {
	return cqrs.NewCommandBusWithConfig(
		publisher,
		cqrs.CommandBusConfig{
			Logger:    logger,
			Marshaler: marshaler,
			GeneratePublishTopic: func(params cqrs.CommandBusGeneratePublishTopicParams) (string, error) {
				return params.CommandName, nil
			},
		},
	)
}

func newCommandProcessor(
	logger watermill.LoggerAdapter,
	marshaler cqrs.CommandEventMarshaler,
	subscriber message.Subscriber,
	router *message.Router,
) (*cqrs.CommandProcessor, error) {
	return cqrs.NewCommandProcessorWithConfig(
		router,
		cqrs.CommandProcessorConfig{
			Logger:    logger,
			Marshaler: marshaler,
			GenerateSubscribeTopic: func(params cqrs.CommandProcessorGenerateSubscribeTopicParams) (string, error) {
				return params.CommandName, nil
			},
			SubscriberConstructor: func(params cqrs.CommandProcessorSubscriberConstructorParams) (message.Subscriber, error) {
				return subscriber, nil
			},
		},
	)
}
