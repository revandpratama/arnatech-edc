package usecase

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/revandpratama/edc-service/internal/dto"
	"github.com/revandpratama/edc-service/internal/entity"
	"github.com/revandpratama/edc-service/internal/repository"
)

type SettlementUsecase interface {
	CreateSettlement(ctx context.Context, saleReq []dto.SaleRequestDTO) (*dto.SettlementResponseDTO, error)
}

type TransactionRepositoryForSettlement interface {
	FindUnsettledByDetails(ctx context.Context, identifiers []repository.TransactionIdentifier) ([]entity.Transaction, error)
}

type settlementUsecase struct {
	settlementRepository               repository.SettlementRepository
	TransactionRepositoryForSettlement TransactionRepositoryForSettlement
}

func NewSettlementUsecase(settlementRepository repository.SettlementRepository, TransactionRepositoryForSettlement TransactionRepositoryForSettlement) SettlementUsecase {
	return &settlementUsecase{
		settlementRepository:               settlementRepository,
		TransactionRepositoryForSettlement: TransactionRepositoryForSettlement,
	}
}

func (s *settlementUsecase) CreateSettlement(ctx context.Context, settlementReq []dto.SaleRequestDTO) (*dto.SettlementResponseDTO, error) {

	var transactionIdentifiers []repository.TransactionIdentifier 
	for _, req := range settlementReq {
		transactionIdentifiers = append(transactionIdentifiers, repository.TransactionIdentifier{
			TerminalID: req.TerminalID,
			Timestamp:  req.Timestamp,
			Amount:     req.Amount,
		})
	}

	transactions, err := s.TransactionRepositoryForSettlement.FindUnsettledByDetails(ctx, transactionIdentifiers)
	if err != nil {
		return nil, err 
	}

	var totalAmount int64
	var approvedCount int
	var declinedCount int
	transactionIDs := []uint{} 
	for _, tx := range transactions {
		if tx.Status == "approved" {
			totalAmount += tx.Amount
			approvedCount++
		} else if tx.Status == "declined" {
			declinedCount++
		}
		transactionIDs = append(transactionIDs, tx.ID)
	}

	batchID := createBatchID()

	newSettlement := &entity.Settlement{
		BatchID:       batchID,
		TotalCount:    len(transactions),
		DeclinedCount: declinedCount,
		ApprovedCount: approvedCount,
		TotalAmount:   totalAmount,
	}

	createdSettlement, err := s.settlementRepository.CreateSettlementWithTransactions(ctx, newSettlement, transactionIDs)
	if err != nil {
		return nil, err
	}

	res := &dto.SettlementResponseDTO{
		BatchID:     createdSettlement.BatchID,
		TotalCount:  createdSettlement.TotalCount,
		Approved:    createdSettlement.ApprovedCount,
		Declined:    createdSettlement.DeclinedCount,
		TotalAmount: createdSettlement.TotalAmount,
	}

	return res, nil
}

func createBatchID() string {
	base := "BATCH"

	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
	if err != nil {
		return ""
	}

	// Combine timestamp (for uniqueness) + random part (for entropy)
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	return fmt.Sprintf("%s%d%06d", base, timestamp, n.Int64())
}
