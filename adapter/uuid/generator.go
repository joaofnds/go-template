package uuid

import (
	"app/internal/id"

	"github.com/google/uuid"
)

var _ id.Generator = (*Generator)(nil)

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (repo *Generator) NewID() string {
	return uuid.NewString()
}
