package casbin

import (
	"app/authz"
	"app/internal/ref"

	"github.com/casbin/casbin/v2"
)

var _ authz.RoleManager = (*RoleManager)(nil)

type RoleManager struct {
	enforcer *casbin.Enforcer
}

func NewRoleManager(enforcer *casbin.Enforcer) *RoleManager {
	return &RoleManager{enforcer: enforcer}
}

func (rolemanager *RoleManager) Assign(user ref.Ref, role ref.Ref) error {
	_, err := rolemanager.enforcer.AddRoleForUserInDomain(user.String(), role.String(), authz.AppDomain)
	return err
}

func (rolemanager *RoleManager) Revoke(user ref.Ref, role ref.Ref) error {
	_, err := rolemanager.enforcer.DeleteRoleForUserInDomain(user.String(), role.String(), authz.AppDomain)
	return err
}

func (rolemanager *RoleManager) GetAll(user ref.Ref) []ref.Ref {
	roleStrings := rolemanager.enforcer.GetRolesForUserInDomain(user.String(), authz.AppDomain)

	roles := make([]ref.Ref, 0, len(roleStrings))
	for _, role := range roleStrings {
		roles = append(roles, ref.NewFromString(role))
	}
	return roles
}

func (rolemanager *RoleManager) RevokeAll(user ref.Ref) error {
	_, err := rolemanager.enforcer.DeleteRolesForUserInDomain(user.String(), authz.AppDomain)
	return err
}
