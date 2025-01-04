package harness

import "go.uber.org/fx"

type Option interface {
	Apply(*Harness)
}

type CallbackOption struct {
	fn func(*Harness)
}

func (option CallbackOption) Apply(harness *Harness) {
	option.fn(harness)
}

func WithCallback(fn func(*Harness)) Option {
	return CallbackOption{fn: fn}
}

func WithoutTX() Option {
	return WithCallback(func(harness *Harness) { harness.useTX = false })
}

func WithoutDeleteAuthUsers() Option {
	return WithCallback(func(harness *Harness) { harness.deleteAuthUsers = false })
}

func WithFxOptions(fxOptions ...fx.Option) Option {
	return WithCallback(func(harness *Harness) {
		harness.fxOptions = append(harness.fxOptions, fxOptions...)
	})
}
