package handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/revandpratama/edc-service/internal/dto"
	"github.com/revandpratama/edc-service/internal/usecase"
)

type SettlementHandler interface {
	CreateSettlement(c *fiber.Ctx) error
}

type settlementHandler struct {
	settlementUsecase usecase.SettlementUsecase
}

func NewSettlementHandler(settlementUsecase usecase.SettlementUsecase) SettlementHandler {
	return &settlementHandler{settlementUsecase: settlementUsecase}
}

func (s *settlementHandler) CreateSettlement(c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	settlementReq, ok := c.Locals("body").(*[]dto.SaleRequestDTO)
	if !ok {
		return c.Status(400).JSON(fiber.Map{"error": "failed to get request body from context"})
	}

	res, err := s.settlementUsecase.CreateSettlement(ctx, *settlementReq)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(res)
}
