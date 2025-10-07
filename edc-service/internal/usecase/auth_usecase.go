package usecase

import (
	"context"

	"github.com/revandpratama/edc-service/internal/dto"
	"github.com/revandpratama/edc-service/internal/repository"
	"github.com/revandpratama/edc-service/util"
)

type AuthUsecase interface {
	GetToken(ctx context.Context, authReq *dto.AuthRequestDTO) (*dto.AuthResponseDTO, error)
}

type authUsecase struct {
	authRepository repository.AuthRepository
}

func NewAuthUsecase(authRepository repository.AuthRepository) AuthUsecase {
	return &authUsecase{
		authRepository: authRepository,
	}
}

func (a *authUsecase) GetToken(ctx context.Context, authReq *dto.AuthRequestDTO) (*dto.AuthResponseDTO, error) {

	err := a.authRepository.CheckAuth(ctx, authReq.TerminalID, authReq.SecretKey)
	if err != nil {
		return nil, err
	}

	// TODO: generate token
	tokenString, err := util.GenerateToken(authReq.TerminalID)
	if err != nil {
		return nil, err
	}

	res := &dto.AuthResponseDTO{
		AccessToken: tokenString,
		TokenType: "Bearer",
	}

	return res, nil
}
