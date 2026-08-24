package watermill

import (
	"app/internal/marshal"
	"uuid"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
)

var _ cqrs.CommandEventMarshaler = Marshaler{}

type Marshaler struct {
	codec marshal.Codec
}

func newMarshaler(codec marshal.Codec) Marshaler {
	return Marshaler{codec: codec}
}

func (m Marshaler) Marshal(v interface{}) (*message.Message, error) {
	b, err := m.codec.Marshal(v)
	if err != nil {
		return nil, err
	}

	msg := message.NewMessage(uuid.New().String(), b)
	msg.Metadata.Set("name", m.Name(v))

	return msg, nil
}

func (m Marshaler) Unmarshal(msg *message.Message, v interface{}) (err error) {
	return m.codec.Unmarshal(msg.Payload, v)
}

func (m Marshaler) Name(cmdOrEvent interface{}) string {
	return cqrs.FullyQualifiedStructName(cmdOrEvent)
}

func (m Marshaler) NameFromMessage(msg *message.Message) string {
	return msg.Metadata.Get("name")
}
