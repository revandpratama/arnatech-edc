package dto

type AuthRequestDTO struct {
	TerminalID string `json:"terminal_id" validate:"required"`
	SecretKey  string `json:"secret_key" validate:"required"`
}

type AuthResponseDTO struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}
