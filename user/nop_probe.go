package user

import (
	"go.uber.org/fx"
)

var NopProbeProvider = fx.Decorate(func() Probe { return NopProbe{} })

type NopProbe struct{}

func (i NopProbe) FailedToCreateUser(_ error)         {}
func (i NopProbe) FailedToDeleteAll(_ error)          {}
func (i NopProbe) FailedToFindByName(_ error)         {}
func (i NopProbe) FailedToRemoveUser(_ error, _ User) {}
func (i NopProbe) UserCreated()                       {}
