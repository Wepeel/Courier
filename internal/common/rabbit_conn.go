package common

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type RabbitConn struct {
	uri  string
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitConn(uri string) (*RabbitConn, error) {
	var rabbitConn RabbitConn = RabbitConn{
		uri: uri,
	}
	var err error
	rabbitConn.conn, err = amqp.Dial(uri)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	rabbitConn.ch, err = rabbitConn.conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
		return nil, err
	}
	return &rabbitConn, nil
}

func (this *RabbitConn) Close() {
	this.ch.Close()
	this.conn.Close()
}

func JsonPrepareObjectForRabbitSend(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}

func ProtobufPrepareObjectForRabbitSend(obj protoreflect.ProtoMessage) ([]byte, error) {
	return proto.Marshal(obj)
}

func (this *RabbitConn) Publish(exchange string, key string, mandatory bool, immediate bool, msg amqp.Publishing) error {
	return this.ch.Publish(exchange, key, mandatory, immediate, msg)
}

func (this *RabbitConn) QueueDeclare(name string, durable bool, autoDelete bool, exclusive bool, noWait bool, args amqp.Table) (amqp.Queue, error) {
	return this.ch.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)
}

func (this *RabbitConn) Consume(queue string, consumer string, autoAck bool, exclusive bool, noLocal bool, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	return this.ch.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
}
