package main

import (
	"log"
	"net"

	pb "github.com/saladin2098/month3/lesson11/public_service/genproto"
	"github.com/saladin2098/month3/lesson11/public_service/service"
	postgres "github.com/saladin2098/month3/lesson11/public_service/storage/postgres"
	"google.golang.org/grpc"
)
func main() {
	db, err := postgres.ConnectDB()
	if err != nil {
		panic(err)
	}
	liss, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCandidateServiceServer(s,service.NewCandidateService(db))
	pb.RegisterElectionServiceServer(s,service.NewElectionService(db))
	pb.RegisterPublicVoteServiceServer(s,service.NewPublicVoteService(db))
	log.Printf("server listening at %v", liss.Addr())
	if err := s.Serve(liss); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
