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

func (p *PromProbe) Listen() {
	event.On(func(e user.UserCreated) { p.UserCreated(e.User) })
	event.On(func(e user.FailedToCreateUser) { p.FailedToCreateUser(e.Err) })
	event.On(func(e user.FailedToDeleteAll) { p.FailedToDeleteAll(e.Err) })
	event.On(func(e user.FailedToFindByName) { p.FailedToFindByName(e.Err) })
	event.On(func(e user.FailedToRemoveUser) { p.FailedToRemoveUser(e.Err, e.User) })
}

func (p *PromProbe) UserCreated(user.User) {
	p.usersCreated.Inc()
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
