package test

import (
	"web/user"

	"go.uber.org/fx"
)

var NopUserInstrumentation = fx.Decorate(NewNopUserInstrumentation)

func NewNopUserInstrumentation() user.Instrumentation {
	return nopUserInstrumentation{}
}

type nopUserInstrumentation struct{}

func (i nopUserInstrumentation) FailedToCreateUser(_ error)              {}
func (i nopUserInstrumentation) FailedToDeleteAll(_ error)               {}
func (i nopUserInstrumentation) FailedToFindByName(_ error)              {}
func (i nopUserInstrumentation) FailedToRemoveUser(_ error, _ user.User) {}
func (l nopUserInstrumentation) UserCreated()                            {}
