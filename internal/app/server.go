package app

import (
	"context"
	"log"
	"math/rand"
	"net"

	pb "github.com/Wepeel/Courier/internal/app/protos"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

const (
	port = ":50051"
)

type Server struct {
	pb.UnimplementedDoctorServiceServer
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func (s *Server) GetDisease(ctx context.Context, in *pb.GetDiseaseRequest) (*pb.GetDiseaseResponse, error) {
	log.Printf("Received %v", in)

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
		return nil, err
	}
	defer ch.Close()

	callback, err := ch.QueueDeclare(
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to create a queue with name %s: %v", callback.Name, err)
		return nil, err
	}

	responses, err := ch.Consume(
		callback.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	corrId := randomString(32)
	in_bytes, err := proto.Marshal(in)
	if err != nil {
		log.Fatalf("Failed to marshal 'in': %v", err)
		return nil, err
	}
	ch.Publish(
		"",
		"rpc_queue",
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       callback.Name,
			Body:          in_bytes,
		})

	for response := range responses {
		if corrId == response.CorrelationId {
			var msg pb.GetDiseaseResponse
			err = proto.Unmarshal(response.Body, &msg)
			if err != nil {
				log.Fatalf("Failed to unmarshal 'response': %v", err)
				return nil, err
			}
			return &msg, nil
		}
	}
	return &pb.GetDiseaseResponse{}, nil
}

func Start() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterDoctorServiceServer(server, &Server{})
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	// TODO: Push request to RabbiqMQ to doctor
	// TODO: doctor Response goes from doctor to hospital
	// TODO: hospital Response goes to courier

}
