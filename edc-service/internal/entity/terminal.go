package entity

import (
	"time"

	"gorm.io/gorm"
)

type Terminal struct {
	ID         uint      `gorm:"primary_key;not null"`
	TerminalID string    `gorm:"type:varchar(100);unique;not null"`
	HMACSecret string    `gorm:"type:varchar(100);not null"`
	CreatedAt  time.Time `gorm:"not null"`
	UpdatedAt  time.Time `gorm:"not null"`
	DeletedAt  gorm.DeletedAt

	MerchantID string `gorm:"not null"`
	Merchant   Merchant
	// Merchant   Merchant `gorm:"foreignKey:MerchantID;references:ID"`
}
