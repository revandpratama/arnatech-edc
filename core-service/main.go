package main

import (
	"fmt"
	"log"
	"net"

	"github.com/revandpratama/core-service/config"
	pb "github.com/revandpratama/core-service/generated/core"
	"github.com/revandpratama/core-service/handler"
	"github.com/revandpratama/core-service/middleware"
	"google.golang.org/grpc"
)

func main() {

	config.LoadConfig()

	GRPC_PORT := fmt.Sprintf(":%s", config.ENV.GRPC_PORT)
	lis, err := net.Listen("tcp", GRPC_PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{
		// Add your token interceptor as a unary interceptor.
		grpc.UnaryInterceptor(middleware.TokenInterceptor),
	}

	s := grpc.NewServer(opts...)

	h := handler.NewHandler()

	pb.RegisterCoreBankingServiceServer(s, h)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

// protoc --proto_path=shared/proto  --go_out=. --go-grpc_out=. shared/proto/core.proto
