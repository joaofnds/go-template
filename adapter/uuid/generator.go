package uuid

import (
	"app/internal/id"
	"uuid"
)

var _ id.Generator = (*Generator)(nil)

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (repo *Generator) NewID() string {
	return uuid.New().String()
}
