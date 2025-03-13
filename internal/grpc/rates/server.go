package rates

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
	"usdt-rate-service/internal/domain"
	"usdt-rate-service/internal/storage/postgres"
	pb "usdt-rate-service/proto/generated"
)

type RateService interface {
	GetRates(context.Context) (domain.Rate, error)
}
type serverAPI struct {
	pb.UnimplementedRateServiceServer
	rateService RateService
	logger      *zap.Logger
	grpc_health_v1.UnimplementedHealthServer
	storage postgres.Storage
}

func Register(gRPCServer *grpc.Server, rateService RateService, logger *zap.Logger, storage postgres.Storage) {
	pb.RegisterRateServiceServer(gRPCServer, &serverAPI{rateService: rateService, logger: logger})
	grpc_health_v1.RegisterHealthServer(gRPCServer, &serverAPI{storage: storage, logger: logger})
}

func (s *serverAPI) GetRates(ctx context.Context, req *pb.GetRatesRequest) (*pb.GetRatesResponse, error) {
	rate, err := s.rateService.GetRates(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get rate")
	}

	return &pb.GetRatesResponse{
		Ask:       rate.Ask,
		Bid:       rate.Bid,
		Timestamp: rate.Timestamp,
	}, nil
}

func (s *serverAPI) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	s.logger.Info("Healthcheck called")

	if !s.isDatabaseAvailable() {
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
		}, nil
	}

	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (s *serverAPI) isDatabaseAvailable() bool {
	err := s.storage.DB.PingContext(context.Background())
	if err != nil {
		s.logger.Error("Database unavailable", zap.Error(err))
		return false
	}
	return true
}
