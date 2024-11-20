package main

import (
	"context"
	"log"
	"net"

	pb "github.com/anishmohan/gin-grpc-service/proto" // Ensure this is the correct import path
	"github.com/gocql/gocql"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCampaignServiceServer
	session *gocql.Session
}

func (s *server) ProcessCampaign(ctx context.Context, req *pb.CampaignRequest) (*pb.CampaignResponse, error) {
	// Log the received data to verify the structure
	log.Printf("Received campaign from user: %s, Campaign Title: %s, Contacts Count: %d",
		req.User.UserName, req.Campaign.Title, len(req.ContactNumbers))

	// Store data into ScyllaDB
	if err := s.session.Query(`INSERT INTO mykeyspace.campaigns (user_id, user_name, user_email, campaign_id, title, description, contact_numbers) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		req.User.UserId, req.User.UserName, req.User.UserEmail, req.Campaign.CampaignId, req.Campaign.Title, req.Campaign.Description, req.ContactNumbers).Exec(); err != nil {
		log.Printf("Failed to insert into ScyllaDB: %v", err)
		return nil, err
	}

	// Simulate processing time of calling external API to send message Twillio or Whatsapp
	// time.Sleep(time.Second * 5)

	response := &pb.CampaignResponse{
		ResponseMessage: "Processed campaign: " + req.Campaign.Title,
	}
	return response, nil
}
func initializeScyllaDB(session *gocql.Session) {
	keyspaceQuery := "CREATE KEYSPACE IF NOT EXISTS mykeyspace WITH replication = {'class': 'NetworkTopologyStrategy', 'AWS_US_EAST_1': 3};"
	if err := session.Query(keyspaceQuery).Exec(); err != nil {
		log.Fatalf("Failed to create keyspace: %v", err)
	}

	tableQuery := `CREATE TABLE IF NOT EXISTS mykeyspace.campaigns (
        user_id text,
        user_name text,
        user_email text,
        campaign_id text PRIMARY KEY,
        title text,
        description text,
        contact_numbers list<text>
    );`
	if err := session.Query(tableQuery).Exec(); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}
func main() {
	// Set up the cluster connection to ScyllaDB
	cluster := gocql.NewCluster(
		"node-0.aws-us-east-1.b99f81492e8076ee62cf.clusters.scylla.cloud",
		"node-1.aws-us-east-1.b99f81492e8076ee62cf.clusters.scylla.cloud",
		"node-2.aws-us-east-1.b99f81492e8076ee62cf.clusters.scylla.cloud",
	)
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: "scylla", Password: "nviY6SIDKJt4z9H"} // Replace with your actual password
	cluster.PoolConfig.HostSelectionPolicy = gocql.DCAwareRoundRobinPolicy("AWS_US_EAST_1")
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to connect to ScyllaDB: %v", err)
	}
	initializeScyllaDB(session)
	defer session.Close()

	// Server setup
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterCampaignServiceServer(grpcServer, &server{session: session})
	log.Println("Server listening at", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
