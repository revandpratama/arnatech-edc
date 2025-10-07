package repository

import (
	"github.com/revandpratama/edc-service/internal/entity"
	"gorm.io/gorm"
)

type MerchantRepository interface {
	GetMerchantByMerchantID(merchantID string) (*entity.Merchant, error)
}

type merchantRepository struct {
	db *gorm.DB
}

func NewMerchantRepository(db *gorm.DB) MerchantRepository {
	return &merchantRepository{db: db}
}

func (r *merchantRepository) GetMerchantByMerchantID(merchantID string) (*entity.Merchant, error) {
	var merchant entity.Merchant
	err := r.db.Where("merchant_id = ?", merchantID).First(&merchant).Error
	return &merchant, err
}
