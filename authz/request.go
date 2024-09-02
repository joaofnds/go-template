package authz

import "app/internal/ref"

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
