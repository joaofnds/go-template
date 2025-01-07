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

func (service *Enforcer) Add(policies ...authz.Policy) error {
	_, err := service.enforcer.AddPoliciesEx(toRules(policies))
	return err
}

func (service *Enforcer) Remove(policies ...authz.Policy) error {
	_, err := service.enforcer.RemovePolicies(toRules(policies))
	return err
}

func toRule(policy authz.Policy) []string {
	return []string{
		policy.Subject.String(),
		policy.Domain,
		policy.Object.String(),
		policy.Action,
		policy.Effect,
	}
}

func toRules(policies []authz.Policy) [][]string {
	rules := make([][]string, len(policies))
	for i, policy := range policies {
		rules[i] = toRule(policy)
	}
	return rules
}
