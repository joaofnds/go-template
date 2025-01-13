package user

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type Emitter struct {
	bus *cqrs.EventBus
}

func NewEmitter(bus *cqrs.EventBus) Emitter {
	return Emitter{bus: bus}
}

func (emitter Emitter) UserCreated(ctx context.Context, user User) error {
	return emitter.bus.Publish(ctx, UserCreated{User: user})
}

func (emitter Emitter) FailedToCreateUser(ctx context.Context, err error) error {
	return emitter.bus.Publish(ctx, FailedToCreateUser{Err: err})
}

func (emitter Emitter) FailedToDeleteAll(ctx context.Context, err error) error {
	return emitter.bus.Publish(ctx, FailedToDeleteAll{Err: err})
}

func (emitter Emitter) FailedToFindByID(ctx context.Context, err error) error {
	return emitter.bus.Publish(ctx, FailedToFindByID{Err: err})
}

func (emitter Emitter) FailedToFindByName(ctx context.Context, err error) error {
	return emitter.bus.Publish(ctx, FailedToFindByName{Err: err})
}

func (emitter Emitter) FailedToRemoveUser(ctx context.Context, err error, user User) error {
	return emitter.bus.Publish(ctx, FailedToRemoveUser{Err: err, User: user})
}

func (emitter Emitter) UserRemoved(ctx context.Context, user User) any {
	return emitter.bus.Publish(ctx, UserRemoved{User: user})
}
