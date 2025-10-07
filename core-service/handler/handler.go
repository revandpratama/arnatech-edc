package handler

import (
	"context"

	pb "github.com/revandpratama/core-service/generated/core"
)

// type Handler interface {
// 	AuthorizeTransaction(ctx context.Context, req *pb.AuthorizeTransactionRequest) (*pb.AuthorizeTransactionResponse, error)
// }

type Handler struct {
	pb.UnimplementedCoreBankingServiceServer
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) AuthorizeTransaction(ctx context.Context, req *pb.AuthorizeTransactionRequest) (*pb.AuthorizeTransactionResponse, error) {

	var status pb.AuthorizationStatus
	var message string

	if req.GetAmount() >= 1000 {
		status = 1
		message = "Transaction approved"
	} else {
		status = 2
		message = "Transaction rejected"
	}

	return &pb.AuthorizeTransactionResponse{
		Status:  status,
		Message: message,
	}, nil
}
