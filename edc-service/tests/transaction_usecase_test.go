package tests

import (
	"context"
	"testing"
	"time"

	pb "github.com/revandpratama/edc-service/generated/core"
	"github.com/revandpratama/edc-service/internal/dto"
	"github.com/revandpratama/edc-service/internal/entity"
	"github.com/revandpratama/edc-service/internal/repository"
	"github.com/revandpratama/edc-service/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)


type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) CreateTransaction(transaction *entity.Transaction) (*entity.Transaction, error) {
	args := m.Called(transaction)
	return args.Get(0).(*entity.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetTransactionByTransactionID(transactionID string) (*entity.Transaction, error) {
	args := m.Called(transactionID)
	return args.Get(0).(*entity.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) ApproveTransaction(ctx context.Context, transaction *entity.Transaction) (*pb.AuthorizeTransactionResponse, error) {
	args := m.Called(ctx, transaction)
	return args.Get(0).(*pb.AuthorizeTransactionResponse), args.Error(1)
}

func (m *MockTransactionRepository) FindUnsettledByDetails(ctx context.Context, identifiers []repository.TransactionIdentifier) ([]entity.Transaction, error) {
	args := m.Called(ctx, identifiers)
	return args.Get(0).([]entity.Transaction), args.Error(1)
}

type MockCoreBankingGateway struct {
	mock.Mock
}

func (m *MockCoreBankingGateway) ApproveTransaction(ctx context.Context, transaction *entity.Transaction) (*pb.AuthorizeTransactionResponse, error) {
	args := m.Called(ctx, transaction)
	return args.Get(0).(*pb.AuthorizeTransactionResponse), args.Error(1)
}

type MockMerchantRepository struct {
	mock.Mock
}

func (m *MockMerchantRepository) GetMerchantByMerchantID(merchantID string) (*entity.Merchant, error) {
	args := m.Called(merchantID)
	return args.Get(0).(*entity.Merchant), args.Error(1)
}

type MockTerminalRepository struct {
	mock.Mock
}

func (m *MockTerminalRepository) GetTerminalByTerminalID(terminalID string) (*entity.Terminal, error) {
	args := m.Called(terminalID)
	return args.Get(0).(*entity.Terminal), args.Error(1)
}

func TestCreateSaleTransaction_Success(t *testing.T) {
	// Arrange: Set up ALL your mocks.
	mockTxnRepo := new(MockTransactionRepository)
	mockMerchantRepo := new(MockMerchantRepository)
	mockTerminalRepo := new(MockTerminalRepository)

	transactionUsecase := usecase.NewTransactionUsecase(mockTxnRepo, mockMerchantRepo, mockTerminalRepo)

	saleReq := &dto.SaleRequestDTO{
		MerchantID: "MCH123",
		TerminalID: "T01",
		Amount:     125000,
		CardNumber: "4111112222221111",
		Timestamp:  time.Now(),
	}

	mockMerchantRepo.On("GetMerchantByMerchantID", saleReq.MerchantID).Return(&entity.Merchant{ID: 1, MerchantID: saleReq.MerchantID}, nil)

	mockTerminalRepo.On("GetTerminalByTerminalID", saleReq.TerminalID).Return(&entity.Terminal{ID: 1, TerminalID: saleReq.TerminalID, MerchantID: "MCH123"}, nil)

	mockTxnRepo.On("ApproveTransaction", mock.Anything, mock.AnythingOfType("*entity.Transaction")).
		Return(&pb.AuthorizeTransactionResponse{Status: pb.AuthorizationStatus_APPROVED}, nil)

	mockTxnRepo.On("CreateTransaction", mock.AnythingOfType("*entity.Transaction")).
		Return(&entity.Transaction{TransactionID: "mock-trx-id"}, nil)

	// Act: Call the method you want to test.
	res, err := transactionUsecase.CreateTransaction(context.Background(), saleReq)

	// Assert: Check the results.
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "approved", res.Status)
	assert.NotEmpty(t, res.TransactionID)

	mockMerchantRepo.AssertExpectations(t)
	mockTerminalRepo.AssertExpectations(t)
	mockTxnRepo.AssertExpectations(t)
}
