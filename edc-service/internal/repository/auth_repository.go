package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/revandpratama/edc-service/internal/entity"
	"gorm.io/gorm"
)

type AuthRepository interface {
	CheckAuth(ctx context.Context, terminalID string, secretKey string) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (a *authRepository) CheckAuth(ctx context.Context, terminalID string, secretKey string) error {
	var terminal entity.Terminal

	err := a.db.WithContext(ctx).Where("terminal_id = ? AND secret_key = ?", terminalID, secretKey).First(&terminal).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("unauthorized: invalid terminal or secret key")
		}
		return err
	}

	return nil
}
