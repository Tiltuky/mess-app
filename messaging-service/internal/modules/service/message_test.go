// messaging-service/internal/modules/service/message_service_test.go
package service

import (
	"errors"
	"messaging-service/internal/models"
	"messaging-service/internal/modules/repository"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type PatternTest struct {
	name          string
	performAction func(*MessageService, *repository.MockChats, *repository.MockMessages)
	verifyResult  func(*testing.T, *MessageService, *repository.MockChats, *repository.MockMessages)
}

func Test_MessageService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChats := repository.NewMockChats(ctrl)
	mockMessages := repository.NewMockMessages(ctrl)

	svc := &MessageService{
		repo:  mockMessages,
		Chats: mockChats,
	}

	tests := []PatternTest{
		{
			name: "SendMessage - success",
			performAction: func(svc *MessageService, chats *repository.MockChats, msgs *repository.MockMessages) {
				chats.EXPECT().GetChatInfo(int64(1)).Return(&models.Chat{}, nil)
				msgs.EXPECT().SendMessage(int64(1), int64(1), "Hello").Return(&models.Message{Content: "Hello"}, nil)
			},
			verifyResult: func(t *testing.T, svc *MessageService, chats *repository.MockChats, msgs *repository.MockMessages) {
				msg, err := svc.SendMessage(1, 1, "Hello")
				assert.NoError(t, err)
				assert.NotNil(t, msg)
				assert.Equal(t, "Hello", msg.Content)
			},
		},
		{
			name: "SendMessage - chat not found",
			performAction: func(svc *MessageService, chats *repository.MockChats, msgs *repository.MockMessages) {
				chats.EXPECT().GetChatInfo(int64(1)).Return(nil, errors.New("not found"))
			},
			verifyResult: func(t *testing.T, svc *MessageService, chats *repository.MockChats, msgs *repository.MockMessages) {
				msg, err := svc.SendMessage(1, 1, "Hello")
				assert.Error(t, err)
				assert.Nil(t, msg)
			},
		},
		{
			name: "GetMessages - success",
			performAction: func(svc *MessageService, chats *repository.MockChats, msgs *repository.MockMessages) {
				chats.EXPECT().GetChatInfo(int64(1)).Return(&models.Chat{}, nil)
				msgs.EXPECT().GetMessages(int64(1)).Return([]models.Message{{Content: "Hello"}}, nil)
			},
			verifyResult: func(t *testing.T, svc *MessageService, chats *repository.MockChats, msgs *repository.MockMessages) {
				messages, err := svc.GetMessages(1)
				assert.NoError(t, err)
				assert.Len(t, messages, 1)
				assert.Equal(t, "Hello", messages[0].Content)
			},
		},
		{
			name: "GetMessages - chat not found",
			performAction: func(svc *MessageService, chats *repository.MockChats, msgs *repository.MockMessages) {
				chats.EXPECT().GetChatInfo(int64(1)).Return(nil, errors.New("not found"))
			},
			verifyResult: func(t *testing.T, svc *MessageService, chats *repository.MockChats, msgs *repository.MockMessages) {
				messages, err := svc.GetMessages(1)
				assert.Error(t, err)
				assert.Nil(t, messages)
			},
		},
		{
			name: "DeleteMessage - success",
			performAction: func(svc *MessageService, chats *repository.MockChats, msgs *repository.MockMessages) {
				chats.EXPECT().GetChatInfo(int64(1)).Return(&models.Chat{}, nil)
				msgs.EXPECT().GetMessages(int64(1)).Return([]models.Message{{ID: 1, Content: "Hello"}}, nil)
				msgs.EXPECT().DeleteMessage(int64(1), int64(1)).Return(nil)
			},
			verifyResult: func(t *testing.T, svc *MessageService, chats *repository.MockChats, msgs *repository.MockMessages) {
				err := svc.DeleteMessage(1, 1)
				assert.NoError(t, err)
			},
		},
		{
			name: "DeleteMessage - message not found",
			performAction: func(svc *MessageService, chats *repository.MockChats, msgs *repository.MockMessages) {
				chats.EXPECT().GetChatInfo(int64(1)).Return(&models.Chat{}, nil)
				msgs.EXPECT().GetMessages(int64(1)).Return(nil, errors.New("not found"))
			},
			verifyResult: func(t *testing.T, svc *MessageService, chats *repository.MockChats, msgs *repository.MockMessages) {
				err := svc.DeleteMessage(1, 1)
				assert.Error(t, err)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.performAction(svc, mockChats, mockMessages)
			test.verifyResult(t, svc, mockChats, mockMessages)
		})
	}
}
