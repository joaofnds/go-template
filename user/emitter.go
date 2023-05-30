package user

import "app/internal/event"

type Emitter struct{}

func NewEventEmitter() Emitter { return Emitter{} }

func (e Emitter) UserCreated(u User) {
	event.Emit(UserCreated{User: u})
}

func (e Emitter) FailedToCreateUser(err error) {
	event.Emit(FailedToCreateUser{Err: err})
}

func (e Emitter) FailedToDeleteAll(err error) {
	event.Emit(FailedToDeleteAll{Err: err})
}

func (e Emitter) FailedToFindByName(err error) {
	event.Emit(FailedToFindByName{Err: err})
}

func (e Emitter) FailedToRemoveUser(err error, u User) {
	event.Emit(FailedToRemoveUser{Err: err, User: u})
}
