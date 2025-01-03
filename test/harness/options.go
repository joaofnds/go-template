package harness

import "go.uber.org/fx"

type Option interface {
	Apply(*Harness)
}

type WithoutTXOption struct{}

func (WithoutTXOption) Apply(harness *Harness) {
	harness.useTX = false
}

func WithoutTX() Option {
	return WithoutTXOption{}
}

type WithoutDeleteAuthUsersOption struct{}

func (WithoutDeleteAuthUsersOption) Apply(harness *Harness) {
	harness.deleteAuthUsers = false
}

func WithoutDeleteAuthUsers() Option {
	return WithoutDeleteAuthUsersOption{}
}

type WithFxOptionsOption struct {
	fxOptions []fx.Option
}

func (option WithFxOptionsOption) Apply(harness *Harness) {
	harness.fxOptions = append(harness.fxOptions, option.fxOptions...)
}

func WithFxOptions(fxOptions ...fx.Option) Option {
	return WithFxOptionsOption{fxOptions: fxOptions}
}
