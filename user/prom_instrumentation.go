package user

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

type PromHabitInstrumentation struct {
	logger            *zap.Logger
	usersCreated      prometheus.Counter
	usersCreateFailed prometheus.Counter
}

func NewPromHabitInstrumentation(logger *zap.Logger) *PromHabitInstrumentation {
	return &PromHabitInstrumentation{
		logger:            logger,
		usersCreated:      promauto.NewCounter(prometheus.CounterOpts{Name: "users_created"}),
		usersCreateFailed: promauto.NewCounter(prometheus.CounterOpts{Name: "users_create_fail"}),
	}
}

func (i *PromHabitInstrumentation) FailedToCreateUser(err error) {
	i.logger.Error("failed to create user", zap.Error(err))
	i.usersCreateFailed.Inc()
}

func (i *PromHabitInstrumentation) FailedToDeleteAll(err error) {
	i.logger.Error("failed to delete all", zap.Error(err))
}

func (i *PromHabitInstrumentation) FailedToFindByName(err error) {
	i.logger.Error("failed to find user by name", zap.Error(err))
}

func (i *PromHabitInstrumentation) FailedToRemoveUser(err error, user User) {
	i.logger.Error("failed to remove user", zap.Error(err), zap.String("name", user.Name))
}

func (l *PromHabitInstrumentation) UserCreated() {
	l.logger.Info("user created")
	l.usersCreated.Inc()
}
