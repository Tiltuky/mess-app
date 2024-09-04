package service

import (
	"messaging-service/internal/models"
	"messaging-service/internal/modules/repository"
)

type ChatsService struct {
	repo repository.Chats
}

func NewChatService(repo repository.Chats) *ChatsService {
	return &ChatsService{repo: repo}
}

func (c *ChatsService) GetChats() (*[]models.Chat, error) {
	return c.repo.GetChats()
}

func (c *ChatsService) CreateChat(chat models.Chat) error {
	return c.repo.CreateChat(chat)
}

func (c *ChatsService) GetChatInfo(chatID int64) (*models.Chat, error) {
	return c.repo.GetChatInfo(chatID)
}

func (c *ChatsService) DeleteChat(chatID int64) error {
	return c.repo.DeleteChat(chatID)
}
