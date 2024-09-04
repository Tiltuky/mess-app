package service

import (
	"errors"
	"messaging-service/internal/models"
	"messaging-service/internal/modules/repository"
)

type MessageService struct {
	repo repository.Messages
	repository.Chats
}

func NewMessageService(repo repository.Messages) *MessageService {
	return &MessageService{repo: repo}
}

func (s *MessageService) SendMessage(chatId, senderId int64, content string) (*models.Message, error) {
	// Проверка существования чата
	chat, err := s.GetChatInfo(chatId)
	if err != nil || chat == nil {
		return nil, errors.New("chat not found")
	}

	// Проверка контента сообщения
	if len(content) == 0 {
		return nil, errors.New("message content cannot be empty")
	}

	return s.repo.SendMessage(chatId, senderId, content)
}

func (s *MessageService) GetMessages(chatId int64) ([]models.Message, error) {
	// Проверка существования чата
	chat, err := s.GetChatInfo(chatId)
	if err != nil || chat == nil {
		return nil, errors.New("chat not found")
	}

	return s.repo.GetMessages(chatId)
}

func (s *MessageService) DeleteMessage(chatID, msgID int64) error {
	// Проверка существования сообщения
	message, err := s.GetMessages(msgID)
	if err != nil || message == nil {
		return errors.New("message not found")
	}

	return s.repo.DeleteMessage(chatID, msgID)
}
