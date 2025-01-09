package casbin

import (
	"app/authz"

	"github.com/casbin/casbin/v2"
)

var _ authz.PermissionManager = (*PermissionManager)(nil)

type PermissionManager struct {
	enforcer *casbin.Enforcer
}

func NewPermissionManager(enforcer *casbin.Enforcer) *PermissionManager {
	return &PermissionManager{enforcer: enforcer}
}

func (permissionManager *PermissionManager) Check(req authz.Request) bool {
	hasPermission, err := permissionManager.enforcer.Enforce(
		req.Subject.String(),
		req.Domain,
		req.Object.String(),
		req.Action,
	)

	return err == nil && hasPermission
}

func (permissionManager *PermissionManager) Add(policies ...authz.Policy) error {
	_, err := permissionManager.enforcer.AddPoliciesEx(toRules(policies))
	return err
}

func (permissionManager *PermissionManager) Remove(policies ...authz.Policy) error {
	_, err := permissionManager.enforcer.RemovePolicies(toRules(policies))
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
