package test

import (
	"web/user"

	"go.uber.org/fx"
)

var Module = fx.Decorate(NewNopUserInstrumentation())

type NopUserInstrumentation struct{}

func NewNopUserInstrumentation() user.Instrumentation {
	return NopUserInstrumentation{}
}

func (i NopUserInstrumentation) FailedToCreateUser(_ error)              {}
func (i NopUserInstrumentation) FailedToDeleteAll(_ error)               {}
func (i NopUserInstrumentation) FailedToFindByName(_ error)              {}
func (i NopUserInstrumentation) FailedToRemoveUser(_ error, _ user.User) {}
func (l NopUserInstrumentation) UserCreated()                            {}
