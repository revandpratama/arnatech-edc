package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/revandpratama/edc-service/config"
	pb "github.com/revandpratama/edc-service/generated/core"
	"github.com/revandpratama/edc-service/internal/adapter"
	"github.com/revandpratama/edc-service/internal/dto"
	"github.com/revandpratama/edc-service/internal/handler"
	"github.com/revandpratama/edc-service/internal/middleware"
	"github.com/revandpratama/edc-service/internal/repository"
	"github.com/revandpratama/edc-service/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	server := NewServer()
	server.Run()

}

type Server struct {
	shutdownCh chan os.Signal
	errorCh    chan error
}

func NewServer() *Server {
	return &Server{
		shutdownCh: make(chan os.Signal, 1),
		errorCh:    make(chan error, 1),
	}
}

func (s *Server) Run() {
	signal.Notify(s.shutdownCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	config.LoadConfig()
	log.Println("Config loaded successfully")

	if err := adapter.ConnectDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Connected to database successfully")

	GRPC_PORT := fmt.Sprintf("core-service:%s", config.ENV.GRPC_PORT)
	conn, err := grpc.NewClient(GRPC_PORT, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	c := pb.NewCoreBankingServiceClient(conn)

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		// Redirects the request to the swagger page
		return c.Redirect("/swagger")
	})

	app.Static("/swagger.yaml", "./api/swagger.yaml")

	app.Get("/swagger/*", swagger.New(swagger.Config{
		// The URL pointing to the yaml file that should be displayed.
		URL: "/swagger.yaml",
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	api := app.Group("/api")
	v1 := api.Group("/v1")

	authRepository := repository.NewAuthRepository(adapter.DB)
	authUsecase := usecase.NewAuthUsecase(authRepository)
	authHandler := handler.NewAuthHandler(authUsecase)
	v1.Post("/auth", middleware.Validate(dto.AuthRequestDTO{}), authHandler.GetToken)

	transactionRepository := repository.NewTransactionRepository(adapter.DB, c)
	merchantRepository := repository.NewMerchantRepository(adapter.DB)
	terminalRepository := repository.NewTerminalRepository(adapter.DB)
	transactionUsecase := usecase.NewTransactionUsecase(transactionRepository, merchantRepository, terminalRepository)

	transactionHandler := handler.NewTransactionHandler(transactionUsecase)
	v1.Post("/transactions/sale", middleware.AuthMiddleware(), middleware.Validate(dto.SaleRequestDTO{}), transactionHandler.CreateTransaction)

	settlementRepository := repository.NewSettlementRepository(adapter.DB)
	settlementUsecase := usecase.NewSettlementUsecase(settlementRepository, transactionRepository)
	settlementHandler := handler.NewSettlementHandler(settlementUsecase)
	v1.Post("/transactions/settlements", middleware.AuthMiddleware(), middleware.Validate([]dto.SaleRequestDTO{}), settlementHandler.CreateSettlement)

	go func() {
		REST_PORT := fmt.Sprintf(":%s", config.ENV.REST_PORT)
		if err := app.Listen(REST_PORT); err != nil {
			s.errorCh <- err
		}

		log.Printf("Server is running on port %s \n", REST_PORT)
	}()

	select {
	case sh := <-s.shutdownCh:
		log.Printf("Shutting down gracefully due to signal: %v \n", sh)
	case err := <-s.errorCh:
		log.Printf("Terminated due to error: %v \n", err)
	}

	if err := app.Shutdown(); err != nil {
		log.Fatalf("Failed to shutdown server: %v", err)
	}

	log.Println("Cleanup complete. Exiting.")
}
