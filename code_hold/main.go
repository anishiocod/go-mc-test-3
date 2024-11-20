package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "github.com/anishmohan/gin-grpc-service/proto" // Ensure this is the correct import path
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCampaignServiceServer
	db *mongo.Collection
}

func (s *server) ProcessCampaign(ctx context.Context, req *pb.CampaignRequest) (*pb.CampaignResponse, error) {
	// Log the received data to verify the structure
	log.Printf("Received campaign from user: %s, Campaign Title: %s, Contacts Count: %d",
		req.User.UserName, req.Campaign.Title, len(req.ContactNumbers))

	// Store data into MongoDB
	_, err := s.db.InsertOne(ctx, bson.M{
		"user": bson.M{
			"userName":  req.User.UserName,
			"userId":    req.User.UserId,
			"userEmail": req.User.UserEmail,
		},
		"campaign": bson.M{
			"campaignId":  req.Campaign.CampaignId,
			"title":       req.Campaign.Title,
			"description": req.Campaign.Description,
		},
		"contactNumbers": req.ContactNumbers,
	})
	if err != nil {
		log.Printf("Failed to insert into MongoDB: %v", err)
		return nil, err
	}

	// Simulate processing time of calling external API to send message Twillio or Whatsapp
	time.Sleep(time.Second * 5)

	response := &pb.CampaignResponse{
		ResponseMessage: "Processed campaign: " + req.Campaign.Title,
	}
	return response, nil
}

func main() {
	// MongoDB client setup
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.TODO())

	// Get handle to MongoDB collection
	collection := client.Database("test").Collection("WA-MS-gRPC")

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterCampaignServiceServer(grpcServer, &server{db: collection})
	log.Println("Server listening at", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
