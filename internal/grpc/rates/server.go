package rates

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"usdt-rate-service/internal/domain"
	pb "usdt-rate-service/proto/generated"
)

type RateService interface {
	GetRates(context.Context) (domain.Rate, error)
}
type serverAPI struct {
	pb.UnimplementedRateServiceServer
	rateService RateService
}

func Register(gRPCServer *grpc.Server, rateService RateService) {
	pb.RegisterRateServiceServer(gRPCServer, &serverAPI{rateService: rateService})
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
