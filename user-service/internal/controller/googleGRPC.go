package controller

import (
	"context"
	pb "user-service/internal/proto/googlePB"
	"user-service/internal/service"
)

type GoogleGRPC interface {
	pb.GoogleServiceServer
	Authenticate(context.Context, *pb.AuthRequest) (*pb.AuthResponse, error)
}

type GoogleGRPCHandler struct {
	pb.UnimplementedGoogleServiceServer
	service service.GoogleServiceInterface
}

func NewGoogleGRPCHandler(service service.GoogleServiceInterface) GoogleGRPC {
	return &GoogleGRPCHandler{service: service}
}

func (s *GoogleGRPCHandler) Authenticate(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	jwt, exp, err := s.service.Authenticate(ctx, req.AccessToken, req.PushToken)
	if err != nil {
		return nil, err
	}

	return &pb.AuthResponse{Token: jwt, ExpiredAt: exp}, nil
}
