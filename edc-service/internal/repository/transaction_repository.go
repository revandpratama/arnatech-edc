package repository

import (
	"context"

	"github.com/revandpratama/edc-service/config"
	pb "github.com/revandpratama/edc-service/generated/core"
	"github.com/revandpratama/edc-service/internal/entity"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateTransaction(transaction *entity.Transaction) (*entity.Transaction, error)
	GetTransactionByTransactionID(transactionID string) (*entity.Transaction, error)
	ApproveTransaction(ctx context.Context, transaction *entity.Transaction) (*pb.AuthorizeTransactionResponse, error)
}

type transactionRepository struct {
	db         *gorm.DB
	grpcClient pb.CoreBankingServiceClient
}

func NewTransactionRepository(db *gorm.DB, grpcClient pb.CoreBankingServiceClient) TransactionRepository {
	return &transactionRepository{
		db:         db,
		grpcClient: grpcClient,
	}
}

func (t *transactionRepository) CreateTransaction(transaction *entity.Transaction) (*entity.Transaction, error) {
	err := t.db.Create(transaction).Error

	return transaction, err
}

func (t *transactionRepository) GetTransactionByTransactionID(transactionID string) (*entity.Transaction, error) {
	var transaction entity.Transaction
	err := t.db.Where("transaction_id = ?", transactionID).First(&transaction).Error
	return &transaction, err
}

func (t *transactionRepository) ApproveTransaction(ctx context.Context, transaction *entity.Transaction) (*pb.AuthorizeTransactionResponse, error) {

	req := &pb.AuthorizeTransactionRequest{
		TransactionId: transaction.TransactionID,
		Amount:        float32(transaction.Amount),
		CardNumber:    transaction.CardNumberMasked,
		Timestamp:     timestamppb.New(transaction.TransactionTimestamp),
		// MerchantId:    transaction.MerchantID,
		TerminalId: transaction.TerminalID,
	}

	token := config.ENV.INTERNAL_GRPC_SECRET
	ctxWithToken := metadata.AppendToOutgoingContext(ctx, "api-token", token)

	response, err := t.grpcClient.AuthorizeTransaction(ctxWithToken, req)

	if err != nil {
		return nil, err
	}

	return response, nil

}
