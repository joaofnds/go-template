package user

import "app/internal/event"

type Emitter struct{}

func NewEventEmitter() Emitter { return Emitter{} }

func (emitter Emitter) UserCreated(user User) {
	event.Emit(UserCreated{User: user})
}

func (emitter Emitter) FailedToCreateUser(err error) {
	event.Emit(FailedToCreateUser{Err: err})
}

func (emitter Emitter) FailedToDeleteAll(err error) {
	event.Emit(FailedToDeleteAll{Err: err})
}

func (emitter Emitter) FailedToFindByName(err error) {
	event.Emit(FailedToFindByName{Err: err})
}

func (emitter Emitter) FailedToRemoveUser(err error, user User) {
	event.Emit(FailedToRemoveUser{Err: err, User: user})
}
