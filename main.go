package main

import (
	"context"
	"log"
	"net"

	pb "github.com/anishmohan/gin-grpc-service/proto" // Update the import path
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCampaignServiceServer
}

func (s *server) ProcessCampaign(ctx context.Context, in *pb.CampaignRequest) (*pb.CampaignResponse, error) {
	log.Printf("Received: %v", in.GetCampaignData())
	// time.Sleep(20 * time.Second) // Simulate some processing time
	return &pb.CampaignResponse{ResponseMessage: "Processed campaign: " + in.GetCampaignData()}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCampaignServiceServer(s, &server{})
	log.Println("server listening at", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
