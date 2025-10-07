package entity

import (
	"time"

	"gorm.io/gorm"
)

type Settlement struct {
	ID            uint      `gorm:"primary_key;not null"`
	BatchID       string    `gorm:"type:varchar(100);unique;not null"`
	TotalCount    int       `gorm:"not null"`
	ApprovedCount int       `gorm:"not null"`
	DeclinedCount int       `gorm:"not null"`
	TotalAmount   int64     `gorm:"type:numeric(15, 2);not null"`
	CreatedAt     time.Time `gorm:"not null"`
	UpdatedAt     time.Time `gorm:"not null"`
	DeletedAt     gorm.DeletedAt

	// Defines a one-to-many relationship
	Transactions []Transaction
}
