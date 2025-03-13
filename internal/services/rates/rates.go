package rates

import (
	"context"
	"go.uber.org/zap"
	"usdt-rate-service/internal/domain"
)

type Rate struct {
	log          *zap.Logger
	rateSaver    RateSaver
	rateProvider RateProvider
}

type RateProvider interface {
	FetchRates(ctx context.Context) (domain.Rate, error)
}

type RateSaver interface {
	SaveRates(ctx context.Context, rate domain.Rate) error
}

func New(log *zap.Logger, rateSaver RateSaver, rateProvider RateProvider) *Rate {
	return &Rate{
		log:          log,
		rateSaver:    rateSaver,
		rateProvider: rateProvider,
	}
}

func (s *Rate) GetRates(ctx context.Context) (domain.Rate, error) {
	s.log.Info("Fetching USDT rates from provider")

	rate, err := s.rateProvider.FetchRates(ctx)
	if err != nil {
		s.log.Error("Failed to fetch USDT rates", zap.Error(err))
		return domain.Rate{}, err
	}

	s.log.Info("Successfully fetched USDT rates from provider")

	s.log.Info("Saving USDT rates to DB")
	err = s.rateSaver.SaveRates(ctx, rate)
	if err != nil {
		s.log.Error("Failed to save USDT rates", zap.Error(err))
		return domain.Rate{}, err
	}

	s.log.Info("Successfully saved USDT rates")

	return rate, nil
}
