package service

import (
	"errors"
	"progekt/dating-app/geolocation-service/internal/models"
	"progekt/dating-app/geolocation-service/internal/modules/geoservice/service/mocks"
	"progekt/dating-app/geolocation-service/internal/modules/geoservice/storage"

	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	h3 "github.com/uber/h3-go/v4"
)

func TestGeoService_UpdateGeolocation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserStorager := mocks.NewMockUserStorager(ctrl)
	mockGeoCacher := mocks.NewMockGeoCacher(ctrl)
	mockMessageQueuer := mocks.NewMockMessageQueuer(ctrl)
	mockPaymentServicer := mocks.NewMockPaymentServicer(ctrl)

	geoService := NewGeoService(mockUserStorager, mockGeoCacher, mockPaymentServicer, mockMessageQueuer, " ", 25*time.Second)

	userID := int64(1)
	latitude := 37.7749
	longitude := -122.4194

	t.Run("успешное обновление", func(t *testing.T) {
		mockGeoCacher.EXPECT().GetUserLocation(userID).Return(nil, storage.ErrUserNotFound)
		mockUserStorager.EXPECT().AddOrUpdateUserLocation(gomock.Any(), gomock.Any()).Return(nil)
		mockGeoCacher.EXPECT().UpdateLocation(gomock.Any(), gomock.Any()).Return(nil)

		message, err := geoService.UpdateGeolocation(userID, latitude, longitude)
		assert.NoError(t, err)
		assert.Equal(t, "Location updated successfully", message)
	})

	t.Run("ошибка при обновлении хранилища", func(t *testing.T) {
		mockGeoCacher.EXPECT().GetUserLocation(userID).Return(nil, storage.ErrUserNotFound)
		mockUserStorager.EXPECT().AddOrUpdateUserLocation(gomock.Any(), gomock.Any()).Return(errors.New("ошибка хранилища"))

		message, err := geoService.UpdateGeolocation(userID, latitude, longitude)
		assert.Error(t, err)
		assert.Equal(t, "", message)
		assert.Contains(t, err.Error(), "failed to update user location in storage")
	})

	t.Run("ошибка при обновлении кэша", func(t *testing.T) {
		mockGeoCacher.EXPECT().GetUserLocation(userID).Return(nil, storage.ErrUserNotFound)
		mockUserStorager.EXPECT().AddOrUpdateUserLocation(gomock.Any(), gomock.Any()).Return(nil)
		mockGeoCacher.EXPECT().UpdateLocation(gomock.Any(), gomock.Any()).Return(errors.New("ошибка кэша"))

		message, err := geoService.UpdateGeolocation(userID, latitude, longitude)
		assert.Error(t, err)
		assert.Equal(t, "", message)
		assert.Contains(t, err.Error(), "failed to update user location in cache")
	})

	t.Run("ошибка при получении местоположения пользователя", func(t *testing.T) {
		mockGeoCacher.EXPECT().GetUserLocation(userID).Return(nil, errors.New("ошибка кэша"))

		message, err := geoService.UpdateGeolocation(userID, latitude, longitude)
		assert.Error(t, err)
		assert.Equal(t, "", message)
		assert.Contains(t, err.Error(), "ошибка кэша")
	})

	t.Run("успешное обновление с существующим пользователем", func(t *testing.T) {
		existingUser := &models.User{
			ID: userID,
			Address: models.Address{
				Lat: latitude,
				Lng: longitude,
			},
			H3Index:   h3.LatLngToCell(h3.LatLng{Lat: latitude, Lng: longitude}, 15).String(),
			CreatedAt: time.Now().Add(-1 * time.Hour),
			UpdatedAt: time.Now().Add(-1 * time.Hour),
		}

		mockGeoCacher.EXPECT().GetUserLocation(userID).Return(existingUser, nil)
		mockUserStorager.EXPECT().AddOrUpdateUserLocation(gomock.Any(), gomock.Any()).Return(nil)
		mockGeoCacher.EXPECT().UpdateLocation(gomock.Any(), gomock.Any()).Return(nil)

		message, err := geoService.UpdateGeolocation(userID, latitude, longitude)
		assert.NoError(t, err)
		assert.Equal(t, "Location updated successfully", message)
	})
}

func TestFindNearbyUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStor := mocks.NewMockUserStorager(ctrl)
	mockCache := mocks.NewMockGeoCacher(ctrl)
	mockPay := mocks.NewMockPaymentServicer(ctrl)
	mockMQ := mocks.NewMockMessageQueuer(ctrl)

	geoService := NewGeoService(mockStor, mockCache, mockPay, mockMQ, "test_topic", 5*time.Second)

	userID := int64(1)
	resolution := int64(9)

	user := &models.User{
		ID: userID,
		Address: models.Address{
			Lat: 37.7749,
			Lng: -122.4194,
		},
	}

	mockCache.EXPECT().GetUserLocation(userID).Return(user, nil)
	mockPay.EXPECT().GetSubscriptionStatus(user).Return("active", nil)

	// Пример ячейки H3
	userCell := h3.LatLngToCell(h3.LatLng{Lat: user.Address.Lat, Lng: user.Address.Lng}, 15)
	parentCell := userCell.Parent(int(resolution))
	childCells := parentCell.Children(userCell.Resolution())

	// Настраиваем ожидания для каждой дочерней ячейки
	for _, cell := range childCells {
		cellString := cell.String()
		mockCache.EXPECT().GetAllUsersInH3Cell(cellString).Return([]int64{2, 3}, nil)
	}

	mockStor.EXPECT().CheckUserIsInActiveSharings(gomock.Any(), userID, gomock.Any()).Return(true, nil).AnyTimes()
	mockCache.EXPECT().GetUserLocation(gomock.Any()).Return(&models.User{
		ID: 2,
		Address: models.Address{
			Lat: 37.7750,
			Lng: -122.4195,
		},
	}, nil).AnyTimes()

	result, err := geoService.FindNearbyUsers(userID, resolution)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.ID)
	assert.NotEmpty(t, result.Nearby)
}

func TestGeoService_GetUserLocation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserStorager := mocks.NewMockUserStorager(ctrl)
	mockGeoCacher := mocks.NewMockGeoCacher(ctrl)
	mockMessageQueuer := mocks.NewMockMessageQueuer(ctrl)
	mockPaymentServicer := mocks.NewMockPaymentServicer(ctrl)

	geoService := NewGeoService(mockUserStorager, mockGeoCacher, mockPaymentServicer, mockMessageQueuer, " ", 25*time.Second)

	responderID := int64(1)
	targetID := int64(2)

	t.Run("successful get user location", func(t *testing.T) {
		mockUserStorager.EXPECT().CheckUserIsInActiveSharings(gomock.Any(), responderID, targetID).Return(true, nil)
		mockGeoCacher.EXPECT().GetUserLocation(targetID).Return(&models.User{ID: targetID}, nil)

		user, err := geoService.GetUserLocation(responderID, targetID)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, targetID, user.ID)
	})

	t.Run("user not in active sharing", func(t *testing.T) {
		mockUserStorager.EXPECT().CheckUserIsInActiveSharings(gomock.Any(), responderID, targetID).Return(false, nil)

		user, err := geoService.GetUserLocation(responderID, targetID)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "user is not in active sharing for you")
	})

	t.Run("error checking active sharing", func(t *testing.T) {
		mockUserStorager.EXPECT().CheckUserIsInActiveSharings(gomock.Any(), responderID, targetID).Return(false, errors.New("storage error"))

		user, err := geoService.GetUserLocation(responderID, targetID)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "failed to check user is in active sharings")
	})
}

func TestGeoService_ShareLocation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserStorager := mocks.NewMockUserStorager(ctrl)
	mockGeoCacher := mocks.NewMockGeoCacher(ctrl)
	mockMessageQueuer := mocks.NewMockMessageQueuer(ctrl)
	mockPaymentServicer := mocks.NewMockPaymentServicer(ctrl)

	geoService := NewGeoService(mockUserStorager, mockGeoCacher, mockPaymentServicer, mockMessageQueuer, " ", 25*time.Second)

	userID := int64(1)
	receiverID := int64(2)
	timeEnd := time.Now().Add(1 * time.Hour)

	t.Run("successful share location", func(t *testing.T) {
		mockUserStorager.EXPECT().ShareLocation(gomock.Any(), userID, gomock.Any()).Return(nil)
		mockMessageQueuer.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)

		err := geoService.ShareLocation(userID, receiverID, timeEnd)
		assert.NoError(t, err)
	})

	t.Run("error sharing location", func(t *testing.T) {
		mockUserStorager.EXPECT().ShareLocation(gomock.Any(), userID, gomock.Any()).Return(errors.New("storage error"))

		err := geoService.ShareLocation(userID, receiverID, timeEnd)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "storage error")
	})
}

func TestGeoService_StopSharingLocation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserStorager := mocks.NewMockUserStorager(ctrl)
	mockGeoCacher := mocks.NewMockGeoCacher(ctrl)
	mockMessageQueuer := mocks.NewMockMessageQueuer(ctrl)
	mockPaymentServicer := mocks.NewMockPaymentServicer(ctrl)

	geoService := NewGeoService(mockUserStorager, mockGeoCacher, mockPaymentServicer, mockMessageQueuer, " ", 25*time.Second)

	userID := int64(1)
	receiverID := int64(2)

	t.Run("successful stop sharing location", func(t *testing.T) {
		mockUserStorager.EXPECT().StopSharingLocation(gomock.Any(), userID, receiverID).Return(nil)

		err := geoService.StopSharingLocation(userID, receiverID)
		assert.NoError(t, err)
	})

	t.Run("error stopping sharing location", func(t *testing.T) {
		mockUserStorager.EXPECT().StopSharingLocation(gomock.Any(), userID, receiverID).Return(errors.New("storage error"))

		err := geoService.StopSharingLocation(userID, receiverID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "storage error")
	})
}

func TestGeoService_SetLocationPrivacy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserStorager := mocks.NewMockUserStorager(ctrl)
	mockGeoCacher := mocks.NewMockGeoCacher(ctrl)
	mockMessageQueuer := mocks.NewMockMessageQueuer(ctrl)
	mockPaymentServicer := mocks.NewMockPaymentServicer(ctrl)

	geoService := NewGeoService(mockUserStorager, mockGeoCacher, mockPaymentServicer, mockMessageQueuer, " ", 25*time.Second)

	userID := int64(1)
	privacy := "private"

	t.Run("successful set location privacy", func(t *testing.T) {
		mockUserStorager.EXPECT().UpdatePrivacy(gomock.Any(), userID, privacy).Return(nil)

		err := geoService.SetLocationPrivacy(userID, privacy)
		assert.NoError(t, err)
	})

	t.Run("error setting location privacy", func(t *testing.T) {
		mockUserStorager.EXPECT().UpdatePrivacy(gomock.Any(), userID, privacy).Return(errors.New("storage error"))

		err := geoService.SetLocationPrivacy(userID, privacy)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "storage error")
	})
}

func TestGeoService_GetLocationHistory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserStorager := mocks.NewMockUserStorager(ctrl)
	mockGeoCacher := mocks.NewMockGeoCacher(ctrl)
	mockMessageQueuer := mocks.NewMockMessageQueuer(ctrl)
	mockPaymentServicer := mocks.NewMockPaymentServicer(ctrl)

	geoService := NewGeoService(mockUserStorager, mockGeoCacher, mockPaymentServicer, mockMessageQueuer, " ", 25*time.Second)

	userID := int64(1)

	t.Run("successful get location history", func(t *testing.T) {
		mockUserStorager.EXPECT().GetLocationHistory(gomock.Any(), userID, 50).Return(&models.User{ID: userID}, nil)

		user, err := geoService.GetLocationHistory(userID)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, userID, user.ID)
	})

	t.Run("error getting location history", func(t *testing.T) {
		mockUserStorager.EXPECT().GetLocationHistory(gomock.Any(), userID, 50).Return(nil, errors.New("storage error"))

		user, err := geoService.GetLocationHistory(userID)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "storage error")
	})
}

func TestGeoService_ClearLocationHistory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserStorager := mocks.NewMockUserStorager(ctrl)
	mockGeoCacher := mocks.NewMockGeoCacher(ctrl)
	mockMessageQueuer := mocks.NewMockMessageQueuer(ctrl)
	mockPaymentServicer := mocks.NewMockPaymentServicer(ctrl)

	geoService := NewGeoService(mockUserStorager, mockGeoCacher, mockPaymentServicer, mockMessageQueuer, " ", 25*time.Second)

	userID := int64(1)

	t.Run("successful clear location history", func(t *testing.T) {
		mockUserStorager.EXPECT().DeleteHistory(gomock.Any(), userID).Return(nil)

		err := geoService.ClearLocationHistory(userID)
		assert.NoError(t, err)
	})

	t.Run("error clearing location history", func(t *testing.T) {
		mockUserStorager.EXPECT().DeleteHistory(gomock.Any(), userID).Return(errors.New("storage error"))

		err := geoService.ClearLocationHistory(userID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "storage error")
	})
}
