package grpcmess

import (
	"context"
	"messaging-service/internal/models"
	proto "messaging-service/proto/messaging"

	"google.golang.org/grpc"
)

type MessServer struct {
	proto.MessagingServiceServer
	messService MessService
}

type MessService interface {
	GetMessages(chatId int64) ([]models.Message, error)
	SendMessage(chatId, senderId int64, content string) (*models.Message, error)
	DeleteMessage(chatID, msgID int64) error
	GetChats() (*[]models.Chat, error)
	CreateChat(chat models.Chat) error
	GetChatInfo(chatID int64) (*models.Chat, error)
	DeleteChat(chatID int64) error
}

func NewMessServer(messService MessService) *MessServer {
	return &MessServer{messService: messService}
}

func Register(gRPC *grpc.Server, mess MessService) {
	proto.RegisterMessagingServiceServer(gRPC, &MessServer{messService: mess})
}

// Получение сообщений в чате по ID чата
func (s *MessServer) GetMessages(ctx context.Context, req *proto.GetMessagesRequest) (*proto.GetMessagesResponse, error) {
	messages, err := s.messService.GetMessages(req.ChatId)
	if err != nil {
		return nil, err
	}

	var protoMessages []*proto.Message
	for _, message := range messages {
		protoMessages = append(protoMessages, &proto.Message{
			Id:        int64(message.ID),
			ChatId:    message.ChatID,
			SenderId:  message.SenderID,
			Content:   message.Content,
			Timestamp: message.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return &proto.GetMessagesResponse{
		Messages: protoMessages,
	}, nil
}

// Отправка сообщения в чат по ID чата
func (s *MessServer) SendMessage(ctx context.Context, req *proto.SendMessageRequest) (*proto.SendMessageResponse, error) {
	message, err := s.messService.SendMessage(req.ChatId, req.SenderId, req.Content)
	if err != nil {
		return nil, err
	}

	return &proto.SendMessageResponse{
		Message: &proto.Message{
			Id:        int64(message.ID),
			ChatId:    message.ChatID,
			SenderId:  message.SenderID,
			Content:   message.Content,
			Timestamp: message.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
		},
	}, nil
}

// Удаление сообщения по ID сообщения в чате
func (s *MessServer) DeleteMessage(ctx context.Context, req *proto.DeleteMessageRequest) (*proto.DeleteMessageResponse, error) {
	err := s.messService.DeleteMessage(req.ChatId, req.MessageId)
	if err != nil {
		return nil, err
	}
	return &proto.DeleteMessageResponse{}, nil
}

// Получение списка чатов текущего пользователя
func (s *MessServer) GetChats(ctx context.Context, req *proto.GetChatsRequest) (*proto.GetChatsResponse, error) {
	chats, err := s.messService.GetChats()
	if err != nil {
		return nil, err
	}

	var protoChats []*proto.Chat
	for _, chat := range *chats {
		protoChats = append(protoChats, &proto.Chat{
			Id:        int64(chat.ID),
			Name:      chat.Name,
			CreatedAt: chat.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return &proto.GetChatsResponse{
		Chats: protoChats,
	}, nil
}

// Создание нового чата
func (s *MessServer) CreateChat(ctx context.Context, req *proto.CreateChatRequest) (*proto.CreateChatResponse, error) {
	chat := models.Chat{
		Name: req.Name,
	}
	err := s.messService.CreateChat(chat)
	if err != nil {
		return nil, err
	}

	return &proto.CreateChatResponse{
		Chat: &proto.Chat{
			Id:        int64(chat.ID),
			Name:      chat.Name,
			CreatedAt: chat.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}, nil
}

// Получение информации о чате по ID
func (s *MessServer) GetChatInfo(ctx context.Context, req *proto.GetChatInfoRequest) (*proto.GetChatInfoResponse, error) {
	chat, err := s.messService.GetChatInfo(req.ChatId)
	if err != nil {
		return nil, err
	}

	return &proto.GetChatInfoResponse{
		Chat: &proto.Chat{
			Id:        int64(chat.ID),
			Name:      chat.Name,
			CreatedAt: chat.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}, nil
}

// Удаление чата по ID
func (s *MessServer) DeleteChat(ctx context.Context, req *proto.DeleteChatRequest) (*proto.DeleteChatResponse, error) {
	err := s.messService.DeleteChat(req.ChatId)
	if err != nil {
		return nil, err
	}
	return &proto.DeleteChatResponse{}, nil
}
