package user

import (
	"go.uber.org/fx"
)

var NopProbeProvider = fx.Decorate(func() Probe { return NopProbe{} })

type NopProbe struct{}

func (p NopProbe) FailedToCreateUser(_ error)         {}
func (p NopProbe) FailedToDeleteAll(_ error)          {}
func (p NopProbe) FailedToFindByName(_ error)         {}
func (p NopProbe) FailedToRemoveUser(_ error, _ User) {}
func (p NopProbe) UserCreated()                       {}
