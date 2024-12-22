package user_adapter

import (
	"app/internal/mill"
	"app/user"
	"context"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
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

func (probe *PromProbe) RegisterEventHandlers(processor *cqrs.EventProcessor) error {
	return processor.AddHandlers(
		mill.NewEventHandler(probe.UserCreated),
		mill.NewEventHandler(probe.FailedToCreateUser),
		mill.NewEventHandler(probe.FailedToDeleteAll),
		mill.NewEventHandler(probe.FailedToFindByName),
		mill.NewEventHandler(probe.FailedToRemoveUser),
	)
}

func (probe *PromProbe) UserCreated(context.Context, *user.UserCreated) error {
	probe.usersCreated.Inc()
	return nil
}

func (probe *PromProbe) FailedToCreateUser(_ context.Context, event *user.FailedToCreateUser) error {
	probe.logger.Error("failed to create user", zap.Error(event.Err))
	probe.usersCreateFailed.Inc()
	return nil
}

func (probe *PromProbe) FailedToDeleteAll(_ context.Context, event *user.FailedToDeleteAll) error {
	probe.logger.Error("failed to delete all", zap.Error(event.Err))
	return nil
}

func (probe *PromProbe) FailedToFindByName(_ context.Context, event *user.FailedToFindByName) error {
	probe.logger.Error("failed to find user by name", zap.Error(event.Err))
	return nil
}

func (probe *PromProbe) FailedToRemoveUser(ctx context.Context, event *user.FailedToRemoveUser) error {
	probe.logger.Error(
		"failed to remove user",
		zap.Error(event.Err),
		zap.String("name", event.User.Email),
	)
	return nil
}
