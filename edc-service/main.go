package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/revandpratama/edc-service/config"
	"github.com/revandpratama/edc-service/internal/adapter"
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
	log.Println(config.ENV.REST_PORT)

	if err := adapter.ConnectDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Connected to database successfully")

	app := fiber.New()

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
