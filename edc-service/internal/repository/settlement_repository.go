package repository

import (
	"context"

	"github.com/revandpratama/edc-service/internal/entity"
	"gorm.io/gorm"
)

type SettlementRepository interface {
	CreateSettlementWithTransactions(ctx context.Context, newSettlement *entity.Settlement, transactionIDs []uint) (*entity.Settlement, error)
}

type settlementRepository struct {
	db *gorm.DB
}

func NewSettlementRepository(db *gorm.DB) SettlementRepository {
	return &settlementRepository{db: db}
}

func (s *settlementRepository) CreateSettlementWithTransactions(ctx context.Context, newSettlement *entity.Settlement, transactionIDs []uint) (*entity.Settlement, error) {

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(newSettlement).Error; err != nil {
			return err
		}

		if err := tx.Model(&entity.Transaction{}).Where("id IN ?", transactionIDs).Update("settlement_id", newSettlement.BatchID).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return newSettlement, err
}
