package main

import (
	"geolocation-service/config"
	"geolocation-service/internal/infrastructure/cache"
	"geolocation-service/internal/infrastructure/db/postgres"
	"geolocation-service/run"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	cfg := config.MustLoad()

	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	db, err := postgres.NewPostgresDB(cfg)
	if err != nil {
		log.Info("Failed to connect to database", zap.Any("err", err))
	}

	redisClient := cache.NewRedisClient(cfg)

	application := run.NewApp(log, cfg, db, redisClient)

	go application.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("shutting down", zap.String("signal", sign.String()))

	application.Stop()
	log.Info("app stopped")
}
