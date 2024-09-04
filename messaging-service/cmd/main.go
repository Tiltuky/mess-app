package main

import (
	"messaging-service/config"
	_ "messaging-service/docs" // Swagger files
	"messaging-service/internal/modules/handler"
	"messaging-service/internal/modules/repository"
	"messaging-service/run"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

// @title Messages and Chats API
// @version 1.0
// @description API для управления сообщениями и чатами

// @host localhost:8080
// @BasePath /

func main() {
	cfg := config.MustLoad()

	log := zap.NewExample()

	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		log.Info("Failed to connect to database", zap.Any("err", err))
	}

	application := run.NewApp(log, cfg, db)

	go application.MustRun()

	go handler.HandleMessages()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("shutting down", zap.String("signal", sign.String()))

	application.Stop()
	log.Info("app stopped")
}
