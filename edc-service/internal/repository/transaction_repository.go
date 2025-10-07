package repository

import (
	"context"
	"fmt"
	"time"

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
	FindUnsettledByDetails(ctx context.Context, identifiers []TransactionIdentifier) ([]entity.Transaction, error)
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

type TransactionIdentifier struct {
	TerminalID string
	Timestamp  time.Time
	Amount     int64
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

func (t *transactionRepository) FindUnsettledByDetails(ctx context.Context, identifiers []TransactionIdentifier) ([]entity.Transaction, error) {
	
	if len(identifiers) == 0 {
		return nil, nil 
	}
	
	var transactions []entity.Transaction

	query := t.db.WithContext(ctx).Where(
        "status IN ? AND settlement_id IS NULL",
        []string{"approved", "declined"},
    )

	andConditions := t.db.Model(&entity.Transaction{})
    for _, id := range identifiers {
        andConditions = andConditions.Or("terminal_id = ? AND transaction_timestamp = ? AND amount = ?", id.TerminalID, id.Timestamp, id.Amount)
    }
    query = query.Where(andConditions)

	if err := query.Find(&transactions).Error; err != nil {
        return nil, err
    }

    if len(transactions) != len(identifiers) {
        return nil, fmt.Errorf("not all transactions found")
    }

    return transactions, nil
	
}
