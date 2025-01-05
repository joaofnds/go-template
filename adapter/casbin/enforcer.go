package casbin

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

func (service *Enforcer) GrantAll(reqs []authz.Request) error {
	casbinRequests := make([][]string, len(reqs))

	for i, req := range reqs {
		casbinRequests[i] = []string{
			req.Subject.String(),
			req.Domain,
			req.Object.String(),
			req.Action,
		}
	}

	_, err := service.enforcer.AddPoliciesEx(casbinRequests)
	return err
}
