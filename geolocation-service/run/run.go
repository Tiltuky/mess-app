package run

import (
	"fmt"

	"geolocation-service/config"
	kafkaConn "geolocation-service/internal/infrastructure/kafka"
	grpcgeo "geolocation-service/internal/modules/geoservice/gRPC"
	"geolocation-service/internal/modules/geoservice/service"
	"geolocation-service/internal/modules/geoservice/storage"
	pServ "geolocation-service/internal/modules/payment/service"
	pStor "geolocation-service/internal/modules/payment/storage"
	"net"

	"github.com/go-redis/redis"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type App struct {
	log        *zap.Logger
	gRPCServer *grpc.Server
	port       int
}

func NewApp(log *zap.Logger, cfg *config.Config, dbPostgres *sqlx.DB, dbRedis *redis.Client) *App {
	kafkaClient, err := kafkaConn.NewKafkaProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	if err != nil {
		log.Fatal("failed to connect kafka", zap.Error(err))
	}
	gRPCServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_zap.UnaryServerInterceptor(log),
		),
		grpc.ChainStreamInterceptor(
			grpc_zap.StreamServerInterceptor(log),
		),
	)
	PostgresStorage := storage.NewUserLocationStorage(dbPostgres)
	RedisStorage := storage.NewGeoCache(dbRedis)
	PaymentStorage := pStor.NewPayStorage(dbPostgres)
	PaymentService := pServ.NewPaymentService(cfg.Payments.SecretKey, cfg.Payments.ProductID, cfg.Payments.TrialEnd, PaymentStorage)

	GeoService := service.NewGeoService(PostgresStorage, RedisStorage, PaymentService, kafkaClient, cfg.Kafka.Topic, cfg.DB.TimeOut)

	grpcgeo.NewGeoServer(GeoService)
	grpcgeo.Register(gRPCServer, GeoService)
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
