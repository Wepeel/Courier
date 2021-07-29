package common

import (
	"encoding/json"

	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type RabbitConn struct {
	uri  string
	conn amqp.Connection
	ch   amqp.Channel
}

func JsonPrepareObjectForRabbitSend(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}

func ProtobufPrepareObjectForRabbitSend(obj protoreflect.ProtoMessage) ([]byte, error) {
	return proto.Marshal(obj)
}

func (this RabbitConn) Send(bytes []byte) {
	return this.ch.Publish()
}
