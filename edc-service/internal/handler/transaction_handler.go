package handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/revandpratama/edc-service/internal/dto"
	"github.com/revandpratama/edc-service/internal/usecase"
)

type TransactionHandler interface {
	CreateTransaction(c *fiber.Ctx) error
}

type transactionHandler struct {
	transactionUsecase usecase.TransactionUsecase
}

func NewTransactionHandler(transactionUsecase usecase.TransactionUsecase) TransactionHandler {
	return &transactionHandler{
		transactionUsecase: transactionUsecase,
	}
}

func (t *transactionHandler) CreateTransaction(c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	saleReq, ok := c.Locals("body").(*dto.SaleRequestDTO)
	if !ok {
		return c.Status(400).JSON(fiber.Map{"error": "failed to get request body from context"})
	}

	res, err := t.transactionUsecase.CreateTransaction(ctx, saleReq)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(res)
}
