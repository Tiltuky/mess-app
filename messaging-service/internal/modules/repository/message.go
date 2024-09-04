package repository

import (
	"messaging-service/internal/models"
	"time"

	"gorm.io/gorm"
)

type MessagePostgres struct {
	db *gorm.DB
}

func NewMessagePostgres(db *gorm.DB) *MessagePostgres {
	return &MessagePostgres{
		db: db,
	}
}

func (r *MessagePostgres) SendMessage(chatId, senderId int64, content string) (*models.Message, error) {
	message := &models.Message{
		ChatID:    chatId,
		SenderID:  senderId,
		Content:   content,
		Timestamp: time.Now(),
	}
	if err := r.db.Create(message).Error; err != nil {
		return nil, err
	}
	return message, nil
}

func (r *MessagePostgres) GetMessages(chatId int64) ([]models.Message, error) {
	var messages []models.Message
	err := r.db.Where("chat_id = ?", chatId).Order("timestamp asc").Find(&messages).Error
	return messages, err
}

func (r *MessagePostgres) DeleteMessage(chatID, msgID int64) error {
	return r.db.Where("chat_id = ? AND id = ?", chatID, msgID).Delete(&models.Message{}).Error
}
