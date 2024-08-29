package time

import (
	"app/internal/clock"
	"time"

	"go.uber.org/fx"
)

var Module = fx.Module("time", fx.Provide(NewClock))

var _ clock.Clock = &Clock{}

type Clock struct{}

func NewClock() *Clock {
	return &Clock{}
}

func (Clock) Now() time.Time {
	return time.Now().UTC()
}
