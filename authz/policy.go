package authz

import "app/internal/ref"

type Policy struct {
	Subject ref.Ref
	Domain  string
	Object  ref.Ref
	Action  string
	Effect  string
}

func NewPolicy(subject ref.Ref, object ref.Ref, action string, effect string) Policy {
	return Policy{
		Subject: subject,
		Domain:  AppDomain,
		Object:  object,
		Action:  action,
		Effect:  effect,
	}
}

func NewAllowPolicy(subject ref.Ref, object ref.Ref, action string) Policy {
	return NewPolicy(subject, object, action, "allow")
}

func NewDenyPolicy(subject ref.Ref, object ref.Ref, action string) Policy {
	return NewPolicy(subject, object, action, "deny")
}
