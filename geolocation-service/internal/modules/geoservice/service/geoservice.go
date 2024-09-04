package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"geolocation-service/internal/models"
	"geolocation-service/internal/modules/geoservice/storage"
	"log"

	"sort"
	"strconv"
	"time"

	"github.com/uber/h3-go/v4"
)

//go:generate mockgen -source=geoservice.go -destination=mocks/mock_geoservice.go -package=mocks

type GeoService struct {
	stor    UserStorager
	cache   GeoCacher
	pay     PaymentServicer
	mq      MessageQueuer
	topic   string
	timeout time.Duration
}

type MessageQueuer interface {
	Publish(topic string, message []byte) error
	Close() error
}

var ErrUserNotFound = errors.New("user not found")

type UserStorager interface {
	AddOrUpdateUserLocation(ctx context.Context, location *models.User) error
	GetUser(ctx context.Context, userID int64) (*models.User, error)
	GetLocationHistory(ctx context.Context, userID int64, limit int) (*models.User, error)
	ShareLocation(ctx context.Context, userID int64, sharing *models.LocationSharing) error
	StopSharingLocation(ctx context.Context, sharerID, receiverID int64) error
	GetActiveSharings(ctx context.Context, userID int64) ([]models.LocationSharing, error)
	UpdatePrivacy(ctx context.Context, userID int64, privacy string) error
	DeleteUser(ctx context.Context, userID int64) error
	DeleteHistory(ctx context.Context, userID int64) error
	CheckUserIsInActiveSharings(ctx context.Context, responder, target int64) (bool, error)
}

type GeoCacher interface {
	UpdateLocation(user models.User, maxAge time.Duration) error
	DeleteUser(userID int64) error
	GetAllUsersInH3Cell(h3Index string) ([]int64, error)
	GetUserLocation(userID int64) (*models.User, error)
	SetPrivacy(ctx context.Context, userID int64, privacy string, maxAge time.Duration) error
	GetUsersWithPrivacy(privacy string) ([]int64, error)
}

type PaymentServicer interface {
	CreateCustomer(user *models.User) (*models.User, error)
	CreateSubscription(user *models.User) error
	GetSubscriptionEndDate(user *models.User) (time.Time, error)
	GetSubscriptionStatus(user *models.User) (string, error)
}

func NewGeoService(stor UserStorager, cache GeoCacher, pay PaymentServicer, mq MessageQueuer, topic string, timeout time.Duration) *GeoService {
	return &GeoService{
		stor:    stor,
		cache:   cache,
		pay:     pay,
		mq:      mq,
		topic:   topic,
		timeout: timeout,
	}
}

func (g *GeoService) UpdateGeolocation(userID int64, latitude, longitude float64) (string, error) {
	// Создаем H3 индекс
	// Результатом этой операции является h3Index - уникальный идентификатор ячейки H3, в которой находится пользователь
	h3Index := h3.LatLngToCell(h3.LatLng{Lat: latitude, Lng: longitude}, 15)

	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()

	var now time.Time
	userLocation, err := g.cache.GetUserLocation(userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			now = time.Now()
		} else {
			return "", err
		}
	} else {
		now = userLocation.CreatedAt
	}

	user := &models.User{
		ID: userID,
		Address: models.Address{
			Lat: latitude,
			Lng: longitude,
		},
		H3Index:   h3Index.String(),
		CreatedAt: now,
		UpdatedAt: time.Now(),
	}

	// Обновляем информацию в хранилище
	err = g.stor.AddOrUpdateUserLocation(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to update user location in storage: %w", err)
	}

	// Обновляем информацию в кэше
	err = g.cache.UpdateLocation(*user, 24*time.Hour) // Предполагаем, что данные актуальны 24 часа
	if err != nil {
		return "", fmt.Errorf("failed to update user location in cache: %w", err)
	}

	return "Location updated successfully", nil
}

func (g *GeoService) FindNearbyUsers(userId int64, resolution int64) (*models.User, error) {
	ctx := context.Background()
	user, err := g.cache.GetUserLocation(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user location: %w", err)
	}

	status, err := g.pay.GetSubscriptionStatus(user)
	if err != nil {
		return nil, fmt.Errorf("failed to get user subscription status: %w", err)
	}

	userCell := h3.LatLngToCell(h3.LatLng{Lat: user.Address.Lat, Lng: user.Address.Lng}, 15)

	// Получаем родительскую ячейку на указанном уровне разрешения
	parentCell := userCell.Parent(int(resolution))

	// Получаем все дочерние ячейки родительской ячейки
	childCells := parentCell.Children(userCell.Resolution())

	user.Nearby = []models.NearbyUser{}

	// Для каждой дочерней ячейки получаем пользователей
	for _, cell := range childCells {
		cellString := cell.String()
		nearbyUserIDs, err := g.cache.GetAllUsersInH3Cell(cellString)
		if err != nil {
			return nil, fmt.Errorf("failed to get users in H3 cell %s: %w", cellString, err)
		}

		for _, nearbyUserID := range nearbyUserIDs {
			// Пропускаем текущего пользователя
			if nearbyUserID == user.ID {
				continue
			}

			shared, err := g.stor.CheckUserIsInActiveSharings(ctx, userId, nearbyUserID)
			if err != nil {
				return nil, fmt.Errorf("failed to check user's sharing status: %w", err)
			}
			if !shared && status != "active" {
				continue
			}

			nearUser, err := g.cache.GetUserLocation(nearbyUserID)
			if err != nil {
				return nil, fmt.Errorf("failed to get user location: %w", err)
			}

			// Преобразуем координаты пользователя в ячейку H3
			nearUserCell := h3.LatLngToCell(h3.LatLng{Lat: nearUser.Address.Lat, Lng: nearUser.Address.Lng}, 15)

			// Рассчитываем расстояние между центрами ячеек
			distance := h3.GreatCircleDistanceM(userCell.LatLng(), nearUserCell.LatLng())

			// Создаем объект NearbyUser и добавляем его в слайс Nearby текущего пользователя
			nearbyUser := models.NearbyUser{
				ID:       nearUser.ID,
				Distance: distance,
			}
			user.Nearby = append(user.Nearby, nearbyUser)
		}
	}

	// Сортируем пользователей по расстоянию
	sort.Slice(user.Nearby, func(i, j int) bool {
		return user.Nearby[i].Distance < user.Nearby[j].Distance
	})

	return user, nil
}

func (g *GeoService) GetUserLocation(responder, target int64) (*models.User, error) {
	ctx := context.Background()
	// Проверяем, что пользователь находится в активной передаче
	isInActiveSharing, err := g.stor.CheckUserIsInActiveSharings(ctx, responder, target)
	if err != nil {
		return nil, fmt.Errorf("failed to check user is in active sharings: %w", err)
	}
	if isInActiveSharing {
		return g.cache.GetUserLocation(target)
	}
	return nil, fmt.Errorf("user is not in active sharing for you")
}

func (g *GeoService) ShareLocation(userId int64, receiverID int64, timeEnd time.Time) error {
	sharing := &models.LocationSharing{
		ReceiverID: receiverID,
		EndTime:    timeEnd,
	}
	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()
	err := g.stor.ShareLocation(ctx, userId, sharing)
	if err != nil {
		return fmt.Errorf("failed to share location: %w", err)
	}

	notification := map[string]interface{}{
		"receiver_id": receiverID,
		"message":     "Location sharing from " + strconv.FormatInt(userId, 10) + " has been started",
	}

	message, err := json.Marshal(notification)
	if err != nil {
		log.Printf("Failed to marshal notification: %v", err)
		return err
	}
	err = g.mq.Publish(g.topic, message)
	if err != nil {
		log.Printf("Failed to publish message to message queue: %v", err)
	}
	return nil
}

func (g *GeoService) StopSharingLocation(userId, receiverID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()
	return g.stor.StopSharingLocation(ctx, userId, receiverID)
}

func (g *GeoService) SetLocationPrivacy(userID int64, privacy string) error {
	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()
	return g.stor.UpdatePrivacy(ctx, userID, privacy)
}

func (g *GeoService) GetLocationHistory(userId int64) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()
	locationHistory, err := g.stor.GetLocationHistory(ctx, userId, 50)
	if err != nil {
		return nil, err
	}
	return locationHistory, nil
}

func (g *GeoService) ClearLocationHistory(userId int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()
	return g.stor.DeleteHistory(ctx, userId)
}
