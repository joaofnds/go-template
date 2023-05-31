package adapter

import (
	"app/internal/event"
	"app/user"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

type PromProbe struct {
	logger            *zap.Logger
	usersCreated      prometheus.Counter
	usersCreateFailed prometheus.Counter
}

func NewPromProbe(logger *zap.Logger) *PromProbe {
	return &PromProbe{
		logger:            logger,
		usersCreated:      promauto.NewCounter(prometheus.CounterOpts{Name: "users_created"}),
		usersCreateFailed: promauto.NewCounter(prometheus.CounterOpts{Name: "users_create_fail"}),
	}
}

func (probe *PromProbe) Listen() {
	event.On(func(event user.UserCreated) { probe.UserCreated(event.User) })
	event.On(func(event user.FailedToCreateUser) { probe.FailedToCreateUser(event.Err) })
	event.On(func(event user.FailedToDeleteAll) { probe.FailedToDeleteAll(event.Err) })
	event.On(func(event user.FailedToFindByName) { probe.FailedToFindByName(event.Err) })
	event.On(func(event user.FailedToRemoveUser) { probe.FailedToRemoveUser(event.Err, event.User) })
}

func (probe *PromProbe) UserCreated(_ user.User) {
	probe.usersCreated.Inc()
}

func (probe *PromProbe) FailedToCreateUser(err error) {
	probe.logger.Error("failed to create user", zap.Error(err))
	probe.usersCreateFailed.Inc()
}

func (probe *PromProbe) FailedToDeleteAll(err error) {
	probe.logger.Error("failed to delete all", zap.Error(err))
}

func (probe *PromProbe) FailedToFindByName(err error) {
	probe.logger.Error("failed to find user by name", zap.Error(err))
}

func (probe *PromProbe) FailedToRemoveUser(err error, failedUser user.User) {
	probe.logger.Error("failed to remove user", zap.Error(err), zap.String("name", failedUser.Name))
}
