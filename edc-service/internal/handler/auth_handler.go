package handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/revandpratama/edc-service/internal/dto"
	"github.com/revandpratama/edc-service/internal/usecase"
)

type AuthHandler interface {
	GetToken(c *fiber.Ctx) error
}

type authHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandler(authUsecase usecase.AuthUsecase) AuthHandler {
	return &authHandler{
		authUsecase: authUsecase,
	}
}

func (a *authHandler) GetToken(c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	authReq, ok := c.Locals("body").(*dto.AuthRequestDTO)
	if !ok {
		return c.Status(400).JSON(fiber.Map{"error": "failed to get request body from context"})
	}

	res, err := a.authUsecase.GetToken(ctx, authReq)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(res)
}
