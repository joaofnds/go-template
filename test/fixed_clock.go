package test

import (
	"app/internal/clock"
	"time"

	"go.uber.org/fx"
)

var Clock = fx.Provide(NewFixedClock)

var _ clock.Clock = FixedClock{}

type FixedClock struct {
	now time.Time
}

func NewFixedClock(now time.Time) FixedClock {
	return FixedClock{now}
}

func (c FixedClock) Now() time.Time {
	return c.now
}
