package watermill

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
)

func newEventBus(
	logger watermill.LoggerAdapter,
	marshaler cqrs.CommandEventMarshaler,
	publisher message.Publisher,
) (*cqrs.EventBus, error) {
	return cqrs.NewEventBusWithConfig(
		publisher,
		cqrs.EventBusConfig{
			Logger:    logger,
			Marshaler: marshaler,
			GeneratePublishTopic: func(params cqrs.GenerateEventPublishTopicParams) (string, error) {
				return params.EventName, nil
			},
		},
	)
}

func newEventProcessor(
	logger watermill.LoggerAdapter,
	marshaler cqrs.CommandEventMarshaler,
	router *message.Router,
	subscriber message.Subscriber,
) (*cqrs.EventProcessor, error) {
	return cqrs.NewEventProcessorWithConfig(
		router,
		cqrs.EventProcessorConfig{
			Logger:    logger,
			Marshaler: marshaler,
			GenerateSubscribeTopic: func(params cqrs.EventProcessorGenerateSubscribeTopicParams) (string, error) {
				return params.EventName, nil
			},
			SubscriberConstructor: func(params cqrs.EventProcessorSubscriberConstructorParams) (message.Subscriber, error) {
				return subscriber, nil
			},
		},
	)
}
