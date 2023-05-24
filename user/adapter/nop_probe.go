package adapter

import (
	"app/user"
	"context"

	"go.uber.org/fx"
)

var NopProbeProvider = fx.Decorate(func() user.Probe { return NopProbe{} })

type NopProbe struct{}

func (p NopProbe) FailedToCreateUser(error)               {}
func (p NopProbe) FailedToDeleteAll(error)                {}
func (p NopProbe) FailedToFindByName(error)               {}
func (p NopProbe) FailedToRemoveUser(error, user.User)    {}
func (p NopProbe) FailedToEnqueue(error)                  {}
func (p NopProbe) UserCreated(context.Context, user.User) {}
