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

type TransactionUsecase interface {
	CreateTransaction(ctx context.Context, saleReq *dto.SaleRequestDTO) (*dto.SaleResponseDTO, error)
}

type TerminalRepositoryForTransaction interface {
	GetTerminalByTerminalID(terminalID string) (*entity.Terminal, error)
}

type MerchantRepositoryForTransaction interface {
	GetMerchantByMerchantID(merchantID string) (*entity.Merchant, error)
}

type transactionUsecase struct {
	transactionRepository            repository.TransactionRepository
	merchantRepositoryForTransaction MerchantRepositoryForTransaction
	terminalRepositoryForTransaction TerminalRepositoryForTransaction
}

func NewTransactionUsecase(transactionRepository repository.TransactionRepository, merchantRepositoryForTransaction MerchantRepositoryForTransaction, terminalRepositoryForTransaction TerminalRepositoryForTransaction) TransactionUsecase {
	return &transactionUsecase{
		transactionRepository:            transactionRepository,
		merchantRepositoryForTransaction: merchantRepositoryForTransaction,
		terminalRepositoryForTransaction: terminalRepositoryForTransaction,
	}
}

func (t *transactionUsecase) CreateTransaction(ctx context.Context, saleReq *dto.SaleRequestDTO) (*dto.SaleResponseDTO, error) {

	// TODO: verify merchant and terminal exists, also if terminal is from that merchant
	merchant, err := t.merchantRepositoryForTransaction.GetMerchantByMerchantID(saleReq.MerchantID)
	if err != nil {
		return nil, err
	}
	terminal, err := t.terminalRepositoryForTransaction.GetTerminalByTerminalID(saleReq.TerminalID)
	if err != nil {
		return nil, err
	}

	if merchant.MerchantID == "" || terminal.MerchantID == "" {
		return nil, fmt.Errorf("merchant or terminal does not exist")
	}

	if merchant.MerchantID != terminal.MerchantID {
		return nil, fmt.Errorf("terminal does not belong to merchant")
	}

	transaction := &entity.Transaction{
		TerminalID:           saleReq.TerminalID,
		Amount:               saleReq.Amount,
		TransactionTimestamp: saleReq.Timestamp,
	}

	//TODO: mask card number after core-service call
	maskedCardNumber, err := maskCardNumber(saleReq.CardNumber)
	if err != nil {
		return nil, err
	}
	transaction.CardNumberMasked = maskedCardNumber

	// TODO: Create transaction ID
	transactionID := createTransactionID()
	transaction.TransactionID = transactionID

	// TODO: Make a call to core-service

	res, err := t.transactionRepository.ApproveTransaction(ctx, transaction)
	if err != nil {
		return nil, err
	}
	var status string
	var message string
	if res.Status == 1 {
		transaction.Status = "approved"
		status = "approved"
		message = "Transaction authorized"
	} else {
		transaction.Status = "rejected"
		status = "rejected"
		message = "Transaction rejected"
	}

	transaction, err = t.transactionRepository.CreateTransaction(transaction)
	if err != nil {
		return nil, err
	}

	saleResponse := dto.SaleResponseDTO{
		TransactionID: transaction.TransactionID,
		Status:        status,
		Message:       message,
	}

	return &saleResponse, nil
}

func maskCardNumber(cardNumber string) (string, error) {
	if len(cardNumber) < 16 {
		return "", fmt.Errorf("invalid card number")
	}

	firstSix := cardNumber[:6]
	lastFour := cardNumber[len(cardNumber)-4:]

	return fmt.Sprintf("%sXXXXXX%s", firstSix, lastFour), nil

}

func createTransactionID() string {
	base := "TRX"

	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
	if err != nil {
		return ""
	}

	// Combine timestamp (for uniqueness) + random part (for entropy)
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	return fmt.Sprintf("%s%d%06d", base, timestamp, n.Int64())
}
