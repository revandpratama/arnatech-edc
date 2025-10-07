package main

import (
	"log"
	"net"

	pb "github.com/revandpratama/core-service/generated/core"
	"github.com/revandpratama/core-service/handler"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	h := handler.NewHandler()

	pb.RegisterCoreBankingServiceServer(s, h)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
// protoc --proto_path=shared/proto  --go_out=. --go-grpc_out=. shared/proto/core.proto
