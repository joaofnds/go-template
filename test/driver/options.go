package driver

import "go.uber.org/fx"

type DriverOption interface {
	Apply(*Driver)
}

type WithoutTXOption struct{}

func (WithoutTXOption) Apply(driver *Driver) {
	driver.useTX = false
}

func WithoutTransaction() DriverOption {
	return WithoutTXOption{}
}

type WithoutDeleteAuthUsersOption struct{}

func (WithoutDeleteAuthUsersOption) Apply(driver *Driver) {
	driver.deleteAuthUsers = false
}

func WithoutDeleteAuthUsers() DriverOption {
	return WithoutDeleteAuthUsersOption{}
}

type WithFxOptionsOption struct {
	fxOptions []fx.Option
}

func (option WithFxOptionsOption) Apply(driver *Driver) {
	driver.fxOptions = append(driver.fxOptions, option.fxOptions...)
}

func WithFxOptions(fxOptions ...fx.Option) DriverOption {
	return WithFxOptionsOption{fxOptions: fxOptions}
}
