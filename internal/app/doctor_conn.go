package app

import (
	"log"

	"github.com/Wepeel/Courier/internal/common"
	"github.com/streadway/amqp"
)

type DoctorConn struct {
	rabbitConn    *common.RabbitConn
	responses     <-chan amqp.Delivery
	callbackQueue string
}

func NewDoctorConn() (*DoctorConn, error) {
	var doctorConn DoctorConn
	var err error
	doctorConn.rabbitConn, err = common.NewRabbitConn("amqp://guest:guest@localhost:5672/")
	callback, err := doctorConn.rabbitConn.QueueDeclare(
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	doctorConn.callbackQueue = callback.Name
	if err != nil {
		log.Fatalf("Failed to create a queue with name %s: %v", callback.Name, err)
		return nil, err
	}

	doctorConn.responses, err = doctorConn.rabbitConn.Consume(
		callback.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	return &doctorConn, err
}

func (this *DoctorConn) Close() {
	this.rabbitConn.Close()
}

func (this *DoctorConn) SendMsgToDoctorConn(msg []byte, corrId string) {
	this.rabbitConn.Publish(
		"",
		"rpc_queue",
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       this.callbackQueue,
			Body:          msg,
		})
}
