package models

import "time"

type Message struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	ChatID    int64     `json:"chat_id" gorm:"not null"`
	SenderID  int64     `json:"sender_id" gorm:"not null"`
	Content   string    `json:"content" gorm:"not null"`
	Timestamp time.Time `json:"timestamp" gorm:"not null"`
}

type Chat struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
}

type ChatMember struct {
	ID     int64 `json:"id" gorm:"primaryKey"`
	ChatID int64 `json:"chat_id" gorm:"not null"`
	UserID int64 `json:"user_id" gorm:"not null"`
}
