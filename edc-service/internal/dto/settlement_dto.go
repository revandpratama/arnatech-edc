package dto

type SettlementRequestDTO struct {
	Transactions []SaleRequestDTO `json:"transactions" validate:"required,min=1"`
}

// SettlementResponseDTO is the response after a settlement batch is processed. [cite: 52]
type SettlementResponseDTO struct {
	BatchID     string  `json:"batch_id"`
	TotalCount  int     `json:"total_count"`
	Approved    int     `json:"approved"`
	Declined    int     `json:"declined"`
	TotalAmount int64 `json:"total_amount"`
}
