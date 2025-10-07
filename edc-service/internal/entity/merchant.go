package entity

import (
	"time"

	"gorm.io/gorm"
)

type Merchant struct {
	ID         uint      `gorm:"primary_key;not null"`
	MerchantID string    `gorm:"type:varchar(50);unique;not null"`
	Name       string    `gorm:"type:varchar(255);not null"`
	CreatedAt  time.Time `gorm:"not null"`
	UpdatedAt  time.Time `gorm:"not null"`
	DeletedAt  gorm.DeletedAt

	Terminals []Terminal
}
