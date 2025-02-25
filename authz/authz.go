package authz

import "app/internal/ref"

const (
	AppDomain = "app"
)

type RoleManager interface {
	Assign(user ref.Ref, role ref.Ref) error
	Revoke(user ref.Ref, role ref.Ref) error
	GetAll(user ref.Ref) []ref.Ref
	RevokeAll(user ref.Ref) error
}

type PermissionManager interface {
	Check(req Request) bool
	Add(policies ...Policy) error
	Remove(policies ...Policy) error
	RemoveBySubject(subject ref.Ref) error
}
