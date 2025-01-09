package user

import (
	"app/authz"
	"app/internal/ref"
)

const (
	PermRead         = "user:read"
	PermDelete       = "user:delete"
	PermReadFeatures = "user:features:read"
)

func NewRef(id string) ref.Ref {
	return ref.New("user", id)
}

type PermissionService struct {
	permissions authz.PermissionManager
}

func NewPermissionService(permissions authz.PermissionManager) *PermissionService {
	return &PermissionService{permissions: permissions}
}

func (service *PermissionService) GrantNewUserPermission(user User) error {
	return service.permissions.Add(
		authz.NewAllowPolicy(NewRef(user.ID), NewRef(user.ID), "*"),
	)
}
