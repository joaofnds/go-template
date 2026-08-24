package watermill

import (
	"uuid"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/bytedance/sonic"
)

var _ cqrs.CommandEventMarshaler = SonicMarshaler{}

type SonicMarshaler struct{}

func newSonicMarshaler() SonicMarshaler {
	return SonicMarshaler{}
}

func (m SonicMarshaler) Marshal(v interface{}) (*message.Message, error) {
	b, err := sonic.Marshal(v)
	if err != nil {
		return nil, err
	}

	msg := message.NewMessage(uuid.New().String(), b)
	msg.Metadata.Set("name", m.Name(v))

	return msg, nil
}

func (SonicMarshaler) Unmarshal(msg *message.Message, v interface{}) (err error) {
	return sonic.Unmarshal(msg.Payload, v)
}

func (m SonicMarshaler) Name(cmdOrEvent interface{}) string {
	return cqrs.FullyQualifiedStructName(cmdOrEvent)
}

func (m SonicMarshaler) NameFromMessage(msg *message.Message) string {
	return msg.Metadata.Get("name")
}
