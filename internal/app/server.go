package app

import (
	"context"
	"log"
	"math/rand"
	"net"

	pb "github.com/Wepeel/Courier/internal/app/protos"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

const (
	port = ":50051"
)

type Server struct {
	pb.UnimplementedDoctorServiceServer
	doctorConn *DoctorConn
}

func (this *Server) Close() {
	this.doctorConn.Close()
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

	corrId := randomString(32)
	in_bytes, err := proto.Marshal(in)
	if err != nil {
		log.Fatalf("Failed to marshal 'in': %v", err)
		return nil, err
	}

	s.doctorConn.SendMsgToDoctorConn(in_bytes, corrId)

	for response := range s.doctorConn.responses {
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

func NewServer() *Server {
	var ret *Server
	var err error
	ret.doctorConn, err = NewDoctorConn()
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
		return nil
	}
	return ret
}

func Start() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	serv := NewServer()
	defer serv.Close()
	pb.RegisterDoctorServiceServer(server, serv)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	// TODO: Push request to RabbiqMQ to doctor
	// TODO: doctor Response goes from doctor to hospital
	// TODO: hospital Response goes to courier

}
