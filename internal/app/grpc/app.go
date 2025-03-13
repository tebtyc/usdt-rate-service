package grpcapp

import (
	"context"
	"fmt"
	"net"

	rategrpc "usdt-rate-service/internal/grpc/rates"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type App struct {
	log        *zap.Logger
	gRPCServer *grpc.Server
	port       string
}

func New(
	log *zap.Logger,
	rateService rategrpc.RateService,
	port string,
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
	))

	rategrpc.Register(gRPCServer, rateService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
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
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", a.port))
	if err != nil {
		return fmt.Errorf("grpcapp.Run: %w", err)
	}

	a.log.Info("grpc server started", zap.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("grpcapp.Run: %w", err)
	}

	return nil
}

func (a *App) Stop() {
	a.log.Info("stopping gRPC server", zap.String("port", a.port))

	a.gRPCServer.GracefulStop()
}
