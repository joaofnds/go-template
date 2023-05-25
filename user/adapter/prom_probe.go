package adapter

import (
	"app/internal/event"
	"app/user"
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

type PromProbe struct {
	logger            *zap.Logger
	userCreatedEvent  *event.Event[user.UserCreated]
	usersCreated      prometheus.Counter
	usersCreateFailed prometheus.Counter
}

func NewPromProbe(logger *zap.Logger, userCreatedEvent *event.Event[user.UserCreated]) *PromProbe {
	return &PromProbe{
		logger:            logger,
		userCreatedEvent:  userCreatedEvent,
		usersCreated:      promauto.NewCounter(prometheus.CounterOpts{Name: "users_created"}),
		usersCreateFailed: promauto.NewCounter(prometheus.CounterOpts{Name: "users_create_fail"}),
	}
}

func (p *PromProbe) FailedToCreateUser(err error) {
	p.logger.Error("failed to create user", zap.Error(err))
	p.usersCreateFailed.Inc()
}

func (p *PromProbe) FailedToDeleteAll(err error) {
	p.logger.Error("failed to delete all", zap.Error(err))
}

func (p *PromProbe) FailedToFindByName(err error) {
	p.logger.Error("failed to find user by name", zap.Error(err))
}

func (p *PromProbe) FailedToRemoveUser(err error, user user.User) {
	p.logger.Error("failed to remove user", zap.Error(err), zap.String("name", user.Name))
}

func (p *PromProbe) FailedToEnqueue(err error) {
	p.logger.Error("failed to enqueue", zap.Error(err))
}

func (p *PromProbe) UserCreated(_ context.Context, u user.User) {
	p.usersCreated.Inc()
	p.userCreatedEvent.Send(user.NewUserCreated(u))
}
