package test

import (
	"github.com/casbin/casbin/v2/persist"
	stringadapter "github.com/casbin/casbin/v2/persist/string-adapter"
	"go.uber.org/fx"
)

var CasbinStringAdapter = fx.Options(
	fx.Decorate(func() persist.Adapter {
		return stringadapter.NewAdapter("g")
	}),
)
