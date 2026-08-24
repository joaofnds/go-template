package json

import (
	"app/internal/marshal"
	"encoding/json/v2"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"json",
	fx.Provide(NewCodec, fx.Private),
	fx.Provide(func(c Codec) marshal.Codec { return c }),
)

var _ marshal.Codec = Codec{}

type Codec struct{}

func NewCodec() Codec {
	return Codec{}
}

func (Codec) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (Codec) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
