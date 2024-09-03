package authzcasbin

import (
	"app/authz"

	"github.com/casbin/casbin/v2"
)

var _ authz.Enforcer = (*Enforcer)(nil)

type Enforcer struct {
	enforcer *casbin.Enforcer
}

func NewEnforcer(enforcer *casbin.Enforcer) *Enforcer {
	return &Enforcer{enforcer: enforcer}
}

func (enforcer *Enforcer) Check(req authz.Request) bool {
	hasPermission, err := enforcer.enforcer.Enforce(
		req.Subject.String(),
		req.Domain,
		req.Object.String(),
		req.Action,
	)

	return err == nil && hasPermission
}

func (service *Enforcer) Grant(req authz.Request) error {
	_, err := service.enforcer.AddPolicy(
		req.Subject.String(),
		req.Domain,
		req.Object.String(),
		req.Action,
	)

	return err
}
