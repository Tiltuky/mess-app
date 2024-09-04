package repository

import (
	"messaging-service/internal/models"

	"gorm.io/gorm"
)

type ChatsPostgres struct {
	db *gorm.DB
}

func NewChatsPostgres(db *gorm.DB) *ChatsPostgres {
	return &ChatsPostgres{
		db: db,
	}

}

func (r *ChatsPostgres) GetChats() (*[]models.Chat, error) {
	var chats []models.Chat
	if err := r.db.Find(&chats).Error; err != nil {
		return nil, err
	}
	return &chats, nil
}

func (r *ChatsPostgres) CreateChat(chat models.Chat) error {
	return r.db.Create(&chat).Error
}

func (r *ChatsPostgres) GetChatInfo(chatId int64) (*models.Chat, error) {
	var chat models.Chat
	if err := r.db.First(&chat, chatId).Error; err != nil {
		return nil, err
	}
	return &chat, nil
}

func (r *ChatsPostgres) DeleteChat(chatId int64) error {
	return r.db.Delete(&models.Chat{}, chatId).Error
}
