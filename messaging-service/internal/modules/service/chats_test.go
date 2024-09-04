// messaging-service/internal/modules/service/chats_service_test.go
package service

import (
	"messaging-service/internal/models"
	"messaging-service/internal/modules/repository"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestChatsService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChatsRepo := repository.NewMockChats(ctrl)
	chatsService := NewChatService(mockChatsRepo)

	t.Run("GetChats - success", func(t *testing.T) {
		expectedChats := &[]models.Chat{{ID: 1, Name: "Test Chat"}}
		mockChatsRepo.EXPECT().GetChats().Return(expectedChats, nil)

		chats, err := chatsService.GetChats()
		assert.NoError(t, err)
		assert.NotNil(t, chats)
		assert.Equal(t, expectedChats, chats)
	})

	t.Run("CreateChat - success", func(t *testing.T) {
		newChat := models.Chat{Name: "New Chat"}
		mockChatsRepo.EXPECT().CreateChat(newChat).Return(nil)

		err := chatsService.CreateChat(newChat)
		assert.NoError(t, err)
	})

	t.Run("GetChatInfo - success", func(t *testing.T) {
		expectedChat := &models.Chat{ID: 1, Name: "Test Chat"}
		mockChatsRepo.EXPECT().GetChatInfo(int64(1)).Return(expectedChat, nil)

		chat, err := chatsService.GetChatInfo(1)
		assert.NoError(t, err)
		assert.NotNil(t, chat)
		assert.Equal(t, expectedChat, chat)
	})

	t.Run("DeleteChat - success", func(t *testing.T) {
		mockChatsRepo.EXPECT().DeleteChat(int64(1)).Return(nil)

		err := chatsService.DeleteChat(1)
		assert.NoError(t, err)
	})
}
