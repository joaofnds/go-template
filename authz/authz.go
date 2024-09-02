package authz

import (
	"app/internal/ref"

	"github.com/casbin/casbin/v2"
)

const (
	AppDomain = "app"
)

type Service struct {
	enforcer *casbin.Enforcer
}

func NewService(enforcer *casbin.Enforcer) *Service {
	return &Service{enforcer: enforcer}
}

func (authz *Service) GrantRole(user ref.Ref, role ref.Ref) error {
	_, err := authz.enforcer.AddRoleForUserInDomain(user.String(), role.String(), AppDomain)
	return err
}

func (authz *Service) RevokeRole(user ref.Ref, role ref.Ref) error {
	_, err := authz.enforcer.DeleteRoleForUserInDomain(user.String(), role.String(), AppDomain)
	return err
}

func (authz *Service) GetRoles(user ref.Ref) []ref.Ref {
	roleStrings := authz.enforcer.GetRolesForUserInDomain(user.String(), AppDomain)

	roles := make([]ref.Ref, 0, len(roleStrings))
	for _, role := range roleStrings {
		roles = append(roles, ref.NewFromString(role))
	}
	return roles
}

func (authz *Service) Check(req Request) bool {
	hasPermission, err := authz.enforcer.Enforce(
		req.Subject.String(),
		req.Domain,
		req.Object.String(),
		req.Action,
	)

	return err == nil && hasPermission
}

func (authz *Service) Grant(req Request) error {
	_, err := authz.enforcer.AddPolicy(
		req.Subject.String(),
		req.Domain,
		req.Object.String(),
		req.Action,
	)

	return err
}
