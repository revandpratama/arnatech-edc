package entity

import (
	"time"

	"gorm.io/gorm"
)

type TransactionStatus string

const (
	StatusPending  TransactionStatus = "pending"
	StatusApproved TransactionStatus = "approved"
	StatusDeclined TransactionStatus = "declined"
)

type Transaction struct {
	ID                   uint              `gorm:"primary_key;not null"`
	TransactionID        string            `gorm:"type:varchar(100);unique;not null"`
	Amount               int64           `gorm:"type:numeric(12, 2);not null"`
	CardNumberMasked     string            `gorm:"type:varchar(20);not null"`
	Status               TransactionStatus `gorm:"type:transaction_status;not null;default:pending"`
	TransactionTimestamp time.Time         `gorm:"not null"`
	CreatedAt            time.Time         `gorm:"not null"`
	UpdatedAt            time.Time         `gorm:"not null"`
	DeletedAt            gorm.DeletedAt

	TerminalID string `gorm:"not null"`
	Terminal   Terminal

	SettlementID *string
	Settlement   *Settlement
}
