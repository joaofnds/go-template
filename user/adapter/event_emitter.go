package adapter

import (
	"app/internal/event"
	"app/user"
)

type EventEmitter struct{}

func NewEventEmitter() *EventEmitter { return &EventEmitter{} }

func (e *EventEmitter) UserCreated(u user.User) {
	event.Send(user.UserCreated{User: u})
}

func (e *EventEmitter) FailedToCreateUser(err error) {
	event.Send(user.FailedToCreateUser{Err: err})
}

func (e *EventEmitter) FailedToDeleteAll(err error) {
	event.Send(user.FailedToDeleteAll{Err: err})
}

func (e *EventEmitter) FailedToFindByName(err error) {
	event.Send(user.FailedToFindByName{Err: err})
}

func (e *EventEmitter) FailedToRemoveUser(err error, u user.User) {
	event.Send(user.FailedToRemoveUser{Err: err, User: u})
}
