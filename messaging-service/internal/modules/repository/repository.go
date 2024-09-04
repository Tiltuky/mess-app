package repository

import (
	"messaging-service/internal/models"

	"gorm.io/gorm"
)

//go:generate mockgen -source=repository.go -destination=mock_repository.go -package=repository
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

type Repository struct {
	Messages
	Chats
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Messages: NewMessagePostgres(db),
		Chats:    NewChatsPostgres(db),
	}
}
