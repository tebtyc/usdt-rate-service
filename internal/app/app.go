package app

import (
	"go.uber.org/zap"
	"usdt-rate-service/internal/lib/exmo"

	grpcapp "usdt-rate-service/internal/app/grpc"
	"usdt-rate-service/internal/services/rates"
	"usdt-rate-service/internal/storage/postgres"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *zap.Logger,
	grpcPort string,
	dbURL string,
) *App {
	storage, err := postgres.New(dbURL)
	if err != nil {
		panic(err)
	}

	provider := exmo.New()

	authService := rates.New(log, storage, provider)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
