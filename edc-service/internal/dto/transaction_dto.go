package dto

import "time"

type SaleRequestDTO struct {
	MerchantID string    `json:"merchant_id" validate:"required"`
	TerminalID string    `json:"terminal_id" validate:"required"`
	Amount     float64   `json:"amount" validate:"required,gt=0"`
	CardNumber string    `json:"card_number" validate:"required"`
	Timestamp  time.Time `json:"timestamp" validate:"required"`
}

type SaleResponseDTO struct {
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
	Message       string `json:"message"`
}
