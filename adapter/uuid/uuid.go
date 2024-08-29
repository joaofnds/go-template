package uuid

import "go.uber.org/fx"

var Module = fx.Module("uuid", fx.Provide(NewGenerator))
