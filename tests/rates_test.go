package tests

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"testing"
	"usdt-rate-service/internal/domain"
	"usdt-rate-service/internal/services/rates"
	"usdt-rate-service/tests/mocks"
)

func TestRate_GetRates_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProvider := mocks.NewMockRateProvider(ctrl)
	mockSaver := mocks.NewMockRateSaver(ctrl)
	logger := zap.NewNop()

	expectedRates := domain.Rate{
		Ask:       1.23,
		Bid:       1.22,
		Timestamp: 1640000000,
	}

	mockProvider.EXPECT().
		FetchRates(gomock.Any()).
		Return(expectedRates, nil)
	mockSaver.EXPECT().
		SaveRates(gomock.Any(), expectedRates).
		Return(nil)

	service := rates.New(logger, mockSaver, mockProvider)

	rate, err := service.GetRates(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, expectedRates, rate)
}

func TestRate_GetRates_FetchError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProvider := mocks.NewMockRateProvider(ctrl)
	mockSaver := mocks.NewMockRateSaver(ctrl)
	logger := zap.NewNop()

	mockProvider.EXPECT().
		FetchRates(gomock.Any()).
		Return(domain.Rate{}, errors.New("fetch error"))

	service := rates.New(logger, mockSaver, mockProvider)
	_, err := service.GetRates(context.Background())

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "fetch error")
}

func TestRate_SaveRates_SaveError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProvider := mocks.NewMockRateProvider(ctrl)
	mockSaver := mocks.NewMockRateSaver(ctrl)
	logger := zap.NewNop()

	expectedRates := domain.Rate{
		Ask:       1.23,
		Bid:       1.22,
		Timestamp: 1640000000,
	}

	mockProvider.EXPECT().
		FetchRates(gomock.Any()).
		Return(expectedRates, nil)
	mockSaver.EXPECT().
		SaveRates(gomock.Any(), expectedRates).
		Return(errors.New("save error"))

	service := rates.New(logger, mockSaver, mockProvider)

	_, err := service.GetRates(context.Background())

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "save error")
}
