package service

import (
	"messaging-service/internal/models"
	"messaging-service/internal/modules/repository"
)

//go:generate mockgen -source=service.go -destination=mock_service.go -package=service
type Messages interface {
	SendMessage(chatId, senderId int64, content string) (*models.Message, error)
	GetMessages(chatId int64) ([]models.Message, error)
	DeleteMessage(chatID, msgID int64) error
}

type Chats interface {
	GetChats() (*[]models.Chat, error)
	CreateChat(chat models.Chat) error
	GetChatInfo(chatID int64) (*models.Chat, error)
	DeleteChat(chatID int64) error
}

type Service struct {
	Messages
	Chats
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Messages: NewMessageService(repos.Messages),
		Chats:    NewChatService(repos.Chats),
	}
}
