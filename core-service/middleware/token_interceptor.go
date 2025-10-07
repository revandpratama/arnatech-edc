package middleware

import (
	"context"

	"github.com/revandpratama/core-service/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TokenInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

	expectedToken := config.ENV.INTERNAL_GRPC_SECRET

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	tokens := md.Get("api-token")
	if len(tokens) < 1 || tokens[0] != expectedToken {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	return handler(ctx, req)
}
