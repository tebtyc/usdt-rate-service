package main

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
	"usdt-rate-service/config"
	"usdt-rate-service/internal/app"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env")
	}

	cfg := config.LoadConfig()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	application := app.New(logger, cfg.GRPC.Port, cfg.GetDBURL())

	go func() {
		application.GRPCServer.MustRun()
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	logger.Info("Shutting down gracefully...")
	application.GRPCServer.Stop()
	logger.Info("Server stopped")
}
