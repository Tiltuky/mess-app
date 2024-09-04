package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user-service/config"
	_ "user-service/docs"
	"user-service/internal/controller"
	"user-service/internal/proto/googlePB"
	"user-service/internal/proto/usersPB"
	"user-service/internal/repository"
	userService "user-service/internal/service"
	_ "user-service/migrations"

	"github.com/gin-gonic/gin"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

const config_file = "config/config.yaml"

type AppInterface interface {
	RunApp()
	AddRoutes()
}

type AppObj struct {
	Config      config.AppConfig
	httpHandler controller.HandlerInterface
	grpcHandler controller.GRPCInterface
	router      *gin.Engine
	grpcServer  *grpc.Server
	logger      *zap.Logger
}

func NewAppObj() *AppObj {
	return &AppObj{httpHandler: controller.NewHandlerObj()}
}

func (a *AppObj) RunApp() {
	// Загружаем конфигурацию приложения
	err := config.LoadConfig(config_file, &a.Config)
	if err != nil {
		log.Fatal(fmt.Errorf("error loading config: %v", err))
	}

	// Создаем errgroup и контекст
	g, ctx := errgroup.WithContext(context.Background())

	// Создаем логгер
	a.logger, err = zap.NewProduction()
	if err != nil {
		log.Fatal("fail to make logger", err)
	}
	defer a.logger.Sync()

	// создаем объект слоя  repository
	repo := repository.NewUserRepoObj(a.logger)

	// создаем объект слоя service
	service := userService.NewUserServiceObj(repo, a.logger)
	gService := userService.NewGoogleService(service, a.logger)

	// создаем объект слоя controller - gRPC хэндлер
	a.grpcHandler = controller.NewGRPCObj(service, a.logger)
	gGRPCHandler := controller.NewGoogleGRPCHandler(gService)

	// запускаем миграции БД
	connString := fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%s sslmode=disable",
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"))
	conn, err := sql.Open("postgres", connString)
	if err != nil {
		a.logger.Fatal("fail open connection to bd", zap.String("error", err.Error()))
	}
	defer conn.Close()
	err = goose.Up(conn, ".")
	if err != nil {
		a.logger.Fatal("fail to migrate db", zap.String("error", err.Error()))
	}

	// создаем и настраиваем http сервер
	addr := fmt.Sprintf("%v:%v", a.Config.HttpServer.HttpHost, a.Config.HttpServer.HttpPort)
	httpServer := http.Server{
		Addr:    addr,
		Handler: a.router,
	}

	// Запускаем http сервер для swagger в отдельной горутине
	g.Go(func() error {
		a.logger.Info(fmt.Sprintf("Server starting on: %v", addr))
		return httpServer.ListenAndServe()
	})

	// создаем и настраиваем gRPC сервер
	a.grpcServer = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_zap.UnaryServerInterceptor(a.logger),
		),
	)

	// регистрируем методы в gRPC сервере
	usersPB.RegisterUsersServiceServer(a.grpcServer, a.grpcHandler)
	googlePB.RegisterGoogleServiceServer(a.grpcServer, gGRPCHandler)

	// Запускаем gRPC сервер в отдельной горутине
	g.Go(func() error {
		connString := fmt.Sprintf("%v:%v", a.Config.GRPCServer.GRPChost, a.Config.GRPCServer.GRPCport)
		a.logger.Info(fmt.Sprintf("Listen on: %v", connString))

		listen, err := net.Listen("tcp", connString)
		if err != nil {
			a.logger.Fatal("error listen :"+connString, zap.String("error", err.Error()))
		}

		return a.grpcServer.Serve(listen)
	})

	// Создаем канал ждущий сигнал для Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	err = httpServer.Shutdown(ctx)
	if err != nil {
		a.logger.Error("Server shutdown error", zap.String("error", err.Error()))
	}
	a.logger.Info("Server stopped gracefully")

}

// Добавляет Http ручки( ручки не реализованны, созданы для сваггера )
func (a *AppObj) AddRoutes() {
	a.router = gin.Default()
	api := a.router.Group("/api")
	users := api.Group("/users")
	{
		users.POST("/register", a.httpHandler.RegisterUser)
		users.POST("/login", a.httpHandler.ListUsers)
		users.GET(":id", a.httpHandler.GetUser)
		users.PUT(":id", a.httpHandler.UpdateUser)
		users.DELETE(":id", a.httpHandler.DeleteUser)
		users.GET(":id/profile", a.httpHandler.GetUserProfile)
		users.PUT(":id/profile", a.httpHandler.UpdateUserProfile)
		users.POST(":id/avatar", a.httpHandler.UploadAvatar)
		users.GET("/", a.httpHandler.LoginUser)
	}

	a.router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
