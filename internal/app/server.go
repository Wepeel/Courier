package app

import (
	"context"
	"log"
	"net"

	pb "github.com/Wepeel/Courier/internal/app/protos"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type Server struct {
	pb.UnimplementedDoctorServiceServer
}

func (s *Server) GetDisease(ctx context.Context, in *pb.GetDiseaseRequest) (*pb.GetDiseaseResponse, error) {
	log.Printf("Received %v", in)
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
	// TODO: hospital Response goes to server
}
