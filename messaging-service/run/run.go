package run

import (
	"fmt"
	"messaging-service/config"
	"messaging-service/internal/modules/repository"
	"messaging-service/internal/modules/service"
	"net"

	grpcmess "messaging-service/internal/modules/gRPC"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type App struct {
	log        *zap.Logger
	gRPCServer *grpc.Server
	port       int
}

func NewApp(log *zap.Logger, cfg *config.Config, dbPostgres *gorm.DB) *App {

	gRPCServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_zap.UnaryServerInterceptor(log),
		),
		grpc.ChainStreamInterceptor(
			grpc_zap.StreamServerInterceptor(log),
		),
	)
	repos := repository.NewRepository(dbPostgres)
	serv := service.NewService(repos)

	grpcmess.NewMessServer(serv)
	grpcmess.Register(gRPCServer, serv)
	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       cfg.Local.Port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.log.Info("grpc server is running", zap.String("address", l.Addr().String()))
	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"
	a.log.With(zap.String("operation", op)).
		Info("grpc server is stopping", zap.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}
