package grpcapp

import (
	"context"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net"
	"net/http"
	"usdt-rate-service/internal/storage/postgres"
	"usdt-rate-service/metrics"

	rategrpc "usdt-rate-service/internal/grpc/rates"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type App struct {
	log              *zap.Logger
	gRPCServer       *grpc.Server
	port             string
	prometheusServer *http.Server
	prometheusPort   string
}

func New(
	log *zap.Logger,
	rateService rategrpc.RateService,
	storage postgres.Storage,
	port string,
	prometheusPort string,
) *App {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.PayloadReceived, logging.PayloadSent,
		),
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			log.Error("Recovered from panic", zap.Any("panic", p))
			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		logging.UnaryServerInterceptor(InterceptorLogger(log), loggingOpts...),
		metrics.UnaryServerInterceptor(),
	))

	metrics.InitMetrics()

	prometheusServer := &http.Server{
		Handler: promhttp.Handler(),
	}

	rategrpc.Register(gRPCServer, rateService, log, storage)

	return &App{
		log:              log,
		gRPCServer:       gRPCServer,
		port:             port,
		prometheusServer: prometheusServer,
		prometheusPort:   prometheusPort,
	}
}

func InterceptorLogger(l *zap.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		zapFields := make([]zap.Field, 0, len(fields))
		for i := 0; i < len(fields); i += 2 {
			if i+1 < len(fields) {
				key, ok := fields[i].(string)
				if !ok {
					key = fmt.Sprintf("unkown_field_%d", i)
				}
				zapFields = append(zapFields, zap.Any(key, fields[i+1]))
			}
		}

		switch lvl {
		case logging.LevelDebug:
			l.Debug(msg, zapFields...)
		case logging.LevelInfo:
			l.Info(msg, zapFields...)
		case logging.LevelWarn:
			l.Warn(msg, zapFields...)
		case logging.LevelError:
			l.Error(msg, zapFields...)
		default:
			l.Info(msg, zapFields...)
		}

	})
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	// gRPC server
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", a.port))
	if err != nil {
		return fmt.Errorf("grpcapp.Run: %w", err)
	}

	a.log.Info("grpc server started", zap.String("addr", l.Addr().String()))

	go func() {
		if err := a.gRPCServer.Serve(l); err != nil {
			a.log.Fatal("grpcapp.Run:", zap.Error(err))
		}
	}()

	// Prometheus metrics server
	a.prometheusServer.Addr = fmt.Sprintf(":%s", a.prometheusPort)

	a.log.Info("prometheus metrics server started", zap.String("addr", a.prometheusServer.Addr))

	go func() {
		if err := a.prometheusServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.log.Fatal("prometheus server error", zap.Error(err))
		}
		a.log.Info("prometheus metrics server stopped", zap.String("addr", a.prometheusServer.Addr))
	}()

	return nil
}

func (a *App) Stop() {
	a.log.Info("stopping Prometheus metrics server", zap.String("port", a.prometheusPort))
	if err := a.prometheusServer.Shutdown(context.Background()); err != nil {
		a.log.Error("failed to shutdown prometheus metrics server", zap.Error(err))
	}

	a.log.Info("stopping gRPC server", zap.String("port", a.port))
	a.gRPCServer.GracefulStop()
}
