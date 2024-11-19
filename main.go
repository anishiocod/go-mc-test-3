package main

import (
	"context"
	"log"
	"net"

	pb "github.com/anishmohan/gin-grpc-service/proto" // ensure this is the correct import path
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCampaignServiceServer
}

func (s *server) ProcessCampaign(ctx context.Context, req *pb.CampaignRequest) (*pb.CampaignResponse, error) {
	// Log the received data to verify the structure
	log.Printf("Received campaign from user: %s, Campaign Title: %s, Contacts Count: %d",
		req.User.UserName, req.Campaign.Title, len(req.ContactNumbers))

	// Simulate processing time
	// time.Sleep(20 * time.Second)

	response := &pb.CampaignResponse{
		ResponseMessage: "Processed campaign: " + req.Campaign.Title,
	}
	return response, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCampaignServiceServer(s, &server{})
	log.Println("Server listening at", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
