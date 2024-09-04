package handler

import (
	"context"
	"log"
	_ "messaging-service/docs"
	responder "messaging-service/internal/modules/handler/responder"
	"messaging-service/internal/modules/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type Router struct {
	*Handler
	*errgroup.Group
	*gin.Engine
	Services *service.Service
}

func NewRouter(log *zap.Logger, e *errgroup.Group, services *service.Service) *Router {
	handler := NewHandler(responder.NewRespond(log))
	return &Router{
		handler,
		e,
		gin.Default(),
		services,
	}
}

func (r *Router) Run(ctx context.Context) {
	messages := r.Engine.Group("/messages")
	{
		messages.GET("/:chat_id", r.getMessages)
		messages.POST("/:chat_id", r.sendMessage)
		messages.DELETE("/:chat_id/:msg_id", r.deleteMessage)
	}

	chats := r.Engine.Group("/chats")
	{
		chats.GET("/", r.getChats)
		chats.POST("/", r.createChat)
		chats.GET("/:id", r.getChatInfo)
		chats.DELETE("/:id", r.deleteChat)
	}

	websocket := r.Engine.Group("/ws")
	{
		websocket.GET("/", r.handleConnection)
	}

	r.Engine.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      r.Engine.Handler(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	sigChan := make(chan os.Signal, 1)
	defer close(sigChan)
	signal.Notify(sigChan, syscall.SIGINT)

	r.Go(func() error {
		log.Println("Starting server...")
		return server.ListenAndServe()
	})

	select {
	case <-sigChan:
		stopCTX, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if err := server.Shutdown(stopCTX); err != nil {
			log.Fatalf("Server shutdown error: %v", err)
		}
	case <-ctx.Done():
		stopCTX, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if err := server.Shutdown(stopCTX); err != nil {
			log.Fatalf("Server shutdown error: %v", err)
		}
	}
}
