package user

import (
	"app/authz"
	"app/internal/ref"
)

func NewRef(id string) ref.Ref {
	return ref.New("user", id)
}

type PermissionService struct {
	enforcer authz.Enforcer
}

func NewPermissionService(enforcer authz.Enforcer) *PermissionService {
	return &PermissionService{enforcer: enforcer}
}

func (service *PermissionService) GrantNewUserPermission(user User) error {
	return service.enforcer.Grant(
		authz.NewAppRequest(NewRef(user.ID), NewRef(user.ID), "*"),
	)
}
