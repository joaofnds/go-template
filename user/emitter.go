package user

import "app/internal/event"

type Emitter struct{}

func NewEventEmitter() Emitter { return Emitter{} }

func (e Emitter) UserCreated(u User) {
	event.Send(UserCreated{User: u})
}

func (e Emitter) FailedToCreateUser(err error) {
	event.Send(FailedToCreateUser{Err: err})
}

func (e Emitter) FailedToDeleteAll(err error) {
	event.Send(FailedToDeleteAll{Err: err})
}

func (e Emitter) FailedToFindByName(err error) {
	event.Send(FailedToFindByName{Err: err})
}

func (e Emitter) FailedToRemoveUser(err error, u User) {
	event.Send(FailedToRemoveUser{Err: err, User: u})
}
