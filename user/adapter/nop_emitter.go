package adapter

import (
	"app/user"

	"go.uber.org/fx"
)

var NopEmitterProvider = fx.Decorate(func() user.Emitter { return NopEmitter{} })

type NopEmitter struct{}

func (e NopEmitter) UserCreated(user.User)               {}
func (e NopEmitter) FailedToCreateUser(error)            {}
func (e NopEmitter) FailedToDeleteAll(error)             {}
func (e NopEmitter) FailedToFindByName(error)            {}
func (e NopEmitter) FailedToRemoveUser(error, user.User) {}
