package user

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type Emitter struct {
	bus *cqrs.EventBus
}

func NewEventEmitter(bus *cqrs.EventBus) Emitter {
	return Emitter{bus: bus}
}

func (emitter Emitter) UserCreated(user User) error {
	return emitter.bus.Publish(context.Background(), UserCreated{User: user})
}

func (emitter Emitter) FailedToCreateUser(err error) error {
	return emitter.bus.Publish(context.Background(), FailedToCreateUser{Err: err})
}

func (emitter Emitter) FailedToDeleteAll(err error) error {
	return emitter.bus.Publish(context.Background(), FailedToDeleteAll{Err: err})
}

func (emitter Emitter) FailedToFindByID(err error) error {
	return emitter.bus.Publish(context.Background(), FailedToFindByID{Err: err})
}

func (emitter Emitter) FailedToFindByName(err error) error {
	return emitter.bus.Publish(context.Background(), FailedToFindByName{Err: err})
}

func (emitter Emitter) FailedToRemoveUser(err error, user User) error {
	return emitter.bus.Publish(context.Background(), FailedToRemoveUser{Err: err, User: user})
}
