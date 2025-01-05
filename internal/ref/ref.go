package ref

import "strings"

type Ref struct {
	ID   string
	Type string
}

func New(id, typ string) Ref {
	return Ref{ID: id, Type: typ}
}

func NewFromString(s string) Ref {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		panic("ref string must have len 2 after split")
	}

	return New(parts[1], parts[0])
}

func (ref Ref) Equal(other Ref) bool {
	return ref.ID == other.ID && ref.Type == other.Type
}

func (ref Ref) String() string {
	return ref.Type + ":" + ref.ID
}
