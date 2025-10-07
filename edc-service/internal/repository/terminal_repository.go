package repository

import (
	"github.com/revandpratama/edc-service/internal/entity"
	"gorm.io/gorm"
)

type TerminalRepository interface {
	GetTerminalByTerminalID(terminalID string) (*entity.Terminal, error)
}

type terminalRepository struct {
	db *gorm.DB
}

func NewTerminalRepository(db *gorm.DB) TerminalRepository {
	return &terminalRepository{db: db}
}

func (t *terminalRepository) GetTerminalByTerminalID(terminalID string) (*entity.Terminal, error) {
	var terminal entity.Terminal
	err := t.db.Where("terminal_id = ?", terminalID).First(&terminal).Error
	return &terminal, err
}
