package user

import (
	"go.uber.org/fx"
)

var NopProbeProvider = fx.Decorate(func() Probe { return NopProbe{} })

type NopProbe struct{}

func (p NopProbe) FailedToCreateUser(error)       {}
func (p NopProbe) FailedToDeleteAll(error)        {}
func (p NopProbe) FailedToFindByName(error)       {}
func (p NopProbe) FailedToRemoveUser(error, User) {}
func (p NopProbe) FailedToEnqueue(error)          {}
func (p NopProbe) UserCreated()                   {}
