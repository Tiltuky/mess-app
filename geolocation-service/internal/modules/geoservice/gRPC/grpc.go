package grpcgeo

import (
	"context"
	"progekt/dating-app/geolocation-service/internal/models"
	"progekt/dating-app/geolocation-service/proto/geolocation/proto"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GeoServer struct {
	proto.GeolocationServiceServer
	geoService GeoService
}

type GeoService interface {
	UpdateGeolocation(userID int64, latitude, longitude float64) (string, error)
	FindNearbyUsers(userId int64, resolution int64) (*models.User, error)
	GetUserLocation(responder, target int64) (*models.User, error)
	ShareLocation(userId int64, receiverID int64, timeEnd time.Time) error
	StopSharingLocation(userId, receiverID int64) error
	SetLocationPrivacy(userID int64, privacy string) error
	GetLocationHistory(userId int64) (*models.User, error)
	ClearLocationHistory(userId int64) error
}

func NewGeoServer(geoService GeoService) *GeoServer {
	return &GeoServer{geoService: geoService}
}

func Register(gRPC *grpc.Server, geo GeoService) {
	proto.RegisterGeolocationServiceServer(gRPC, &GeoServer{geoService: geo})
}

// Обновление геолокации текущего пользователя
func (s *GeoServer) UpdateGeolocation(ctx context.Context, req *proto.UpdateGeolocationRequest) (*proto.UpdateGeolocationResponse, error) {
	message, err := s.geoService.UpdateGeolocation(req.UserId, req.Latitude, req.Longitude)
	if err != nil {
		return nil, err
	}
	return &proto.UpdateGeolocationResponse{Message: message}, nil
}

// Поиск пользователей поблизости
func (s *GeoServer) FindNearbyUsers(ctx context.Context, req *proto.FindNearbyUsersRequest) (*proto.FindNearbyUsersResponse, error) {
	users, err := s.geoService.FindNearbyUsers(req.UserId, req.Resolution)
	if err != nil {
		return nil, err
	}

	response := &proto.FindNearbyUsersResponse{}
	for _, user := range users.Nearby {
		response.Users = append(response.Users, &proto.NearbyUser{
			Id:       user.ID,
			Distance: user.Distance,
		})
	}
	return response, nil
}

// Получение текущей геолокации пользователя по его ID
func (s *GeoServer) GetUserLocation(ctx context.Context, req *proto.GetUserLocationRequest) (*proto.GetUserLocationResponse, error) {
	location, err := s.geoService.GetUserLocation(req.Idresp, req.IdTarget)
	if err != nil {
		return nil, err
	}
	return &proto.GetUserLocationResponse{
		Location: &proto.UserLocation{
			Id:        location.ID,
			Latitude:  location.Address.Lat,
			Longitude: location.Address.Lng,
		},
	}, nil
}

// Поделиться своей геолокацией с другим пользователем
func (s *GeoServer) ShareLocation(ctx context.Context, req *proto.ShareLocationRequest) (*proto.ShareLocationResponse, error) {
	err := s.geoService.ShareLocation(req.Idresp, req.IdTarget, req.TimeEnd.AsTime())
	if err != nil {
		return nil, err
	}
	message := "Геолокация успешно расшарена"
	return &proto.ShareLocationResponse{Message: message}, nil
}

// Прекратить делиться своей геолокацией с другим пользователем
func (s *GeoServer) StopSharingLocation(ctx context.Context, req *proto.StopSharingLocationRequest) (*proto.StopSharingLocationResponse, error) {
	err := s.geoService.StopSharingLocation(req.UserId, req.ReceiverId)
	if err != nil {
		return nil, err
	}
	message := "Геолокация успешно остановлена"
	return &proto.StopSharingLocationResponse{Message: message}, nil
}

// Настройка конфиденциальности геолокации
func (s *GeoServer) SetLocationPrivacy(ctx context.Context, req *proto.SetLocationPrivacyRequest) (*proto.SetLocationPrivacyResponse, error) {
	err := s.geoService.SetLocationPrivacy(req.UserId, req.Visibility)
	if err != nil {
		return nil, err
	}
	message := "Конфиденциальность геолокации успешно изменена"
	return &proto.SetLocationPrivacyResponse{Message: message}, nil
}

// Получение истории геолокации пользователя
func (s *GeoServer) GetLocationHistory(ctx context.Context, req *proto.GetLocationHistoryRequest) (*proto.GetLocationHistoryResponse, error) {
	history, err := s.geoService.GetLocationHistory(req.UserId)
	if err != nil {
		return nil, err
	}

	response := &proto.GetLocationHistoryResponse{}
	for _, entry := range history.LHistory {
		response.History = append(response.History, &proto.LocationHistory{
			Timestamp: timestamppb.New(entry.Timestamp),
			Latitude:  entry.Address.Lat,
			Longitude: entry.Address.Lng,
		})
	}
	return response, nil
}

// Очистка истории геолокации пользователя
func (s *GeoServer) ClearLocationHistory(ctx context.Context, req *proto.ClearLocationHistoryRequest) (*proto.ClearLocationHistoryResponse, error) {
	err := s.geoService.ClearLocationHistory(req.UserId)
	if err != nil {
		return nil, err
	}
	message := "История геолокации пользователя успешно очищена"
	return &proto.ClearLocationHistoryResponse{Message: message}, nil
}
