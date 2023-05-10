package user

import (
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

func (i *PromProbe) FailedToCreateUser(err error) {
	i.logger.Error("failed to create user", zap.Error(err))
	i.usersCreateFailed.Inc()
}

func (i *PromProbe) FailedToDeleteAll(err error) {
	i.logger.Error("failed to delete all", zap.Error(err))
}

func (i *PromProbe) FailedToFindByName(err error) {
	i.logger.Error("failed to find user by name", zap.Error(err))
}

func (i *PromProbe) FailedToRemoveUser(err error, user User) {
	i.logger.Error("failed to remove user", zap.Error(err), zap.String("name", user.Name))
}

func (l *PromProbe) UserCreated() {
	l.logger.Info("user created")
	l.usersCreated.Inc()
}
