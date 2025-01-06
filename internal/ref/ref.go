package ref

import "strings"

type Ref struct {
	Type string
	ID   string
}

func New(typ, id string) Ref {
	return Ref{Type: typ, ID: id}
}

func NewFromString(s string) Ref {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		panic("ref string must have len 2 after split")
	}
	if parts[0] == "" || parts[1] == "" {
		panic("ref string must have non-empty parts")
	}

	return New(parts[0], parts[1])
}

func (ref Ref) Equal(other Ref) bool {
	return ref.Type == other.Type && ref.ID == other.ID
}

func (ref Ref) String() string {
	return ref.Type + ":" + ref.ID
}
