package metrics

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"time"
)

// Определение метрик
var (
	grpcRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_requests_total",
			Help: "Total number of gRPC requests",
		},
		[]string{"method", "status"},
	)

	grpcDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_duration_seconds",
			Help:    "Histogram of gRPC request durations",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)
)

// Инициализация метрик
func InitMetrics() {
	prometheus.MustRegister(grpcRequests)
	prometheus.MustRegister(grpcDuration)
}

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		// Выполняем обработку запроса
		resp, err := handler(ctx, req)

		// Время обработки запроса
		duration := time.Since(start).Seconds()

		// Запись метрик
		grpcRequests.WithLabelValues(info.FullMethod, status.Code(err).String()).Inc()
		grpcDuration.WithLabelValues(info.FullMethod).Observe(duration)

		return resp, err
	}
}
