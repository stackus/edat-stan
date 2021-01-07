package edatstan

import (
	"github.com/gogo/protobuf/proto"
	"github.com/nats-io/stan.go"

	"github.com/stackus/edat/msg"
)

var DefaultSerializer = ProtoSerializer{}

type Serializer interface {
	Serialize(message msg.Message) ([]byte, error)
	Deserialize(message *stan.Msg) (msg.Message, error)
}

type ProtoSerializer struct{}

func (ProtoSerializer) Serialize(message msg.Message) ([]byte, error) {
	return proto.Marshal(&Msg{
		Id:      message.ID(),
		Headers: message.Headers(),
		Payload: message.Payload(),
	})
}

func (ProtoSerializer) Deserialize(message *stan.Msg) (msg.Message, error) {
	protoMsg := &Msg{}

	err := proto.Unmarshal(message.Data, protoMsg)
	if err != nil {
		return nil, err
	}

	return msg.NewMessage(protoMsg.Payload, msg.WithMessageID(protoMsg.Id), msg.WithHeaders(protoMsg.Headers)), nil
}
