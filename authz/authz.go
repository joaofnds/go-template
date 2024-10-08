package authz

import "app/internal/ref"

const (
	AppDomain = "app"
)

type RoleManager interface {
	Assign(user ref.Ref, role ref.Ref) error
	Revoke(user ref.Ref, role ref.Ref) error
	GetAll(user ref.Ref) []ref.Ref
}

type Enforcer interface {
	Grant(req Request) error
	Check(req Request) bool
}

type Request struct {
	Subject ref.Ref
	Domain  string
	Object  ref.Ref
	Action  string
}

func NewAppRequest(subject ref.Ref, object ref.Ref, action string) Request {
	return Request{
		Subject: subject,
		Domain:  AppDomain,
		Object:  object,
		Action:  action,
	}
}
