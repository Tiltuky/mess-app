package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"geolocation-service/internal/models"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

var ErrUserNotFound = errors.New("user not found")

type GeoCache struct {
	rdb *redis.Client
}

func NewGeoCache(redisClient *redis.Client) *GeoCache {
	return &GeoCache{
		rdb: redisClient,
	}
}

// UpdateLocation обновляет местоположение пользователя в Redis
func (c *GeoCache) UpdateLocation(user models.User, maxAge time.Duration) error {
	if user.ID <= 0 {
		return fmt.Errorf("invalid user ID: %d", user.ID)
	}
	if user.H3Index == "" || len(user.H3Index) != 15 {
		return fmt.Errorf("invalid H3Index: %s", user.H3Index)
	}
	// Получение значения по ключу из геоиндекса

	pipe := c.rdb.Pipeline()

	// Удаляем пользователя из геоиндекса
	pipe.ZRem(fmt.Sprintf("h3:%s", user.H3Index), user.ID)
	// Обновляем геоиндекс
	pipe.GeoAdd(fmt.Sprintf("h3:%s", user.H3Index), &redis.GeoLocation{
		Name:      fmt.Sprintf("%d", user.ID),
		Longitude: user.Address.Lng,
		Latitude:  user.Address.Lat,
	})

	// Сохраняем полную информацию о местоположении пользователя
	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}
	pipe.Set(fmt.Sprintf("user_location:%d", user.ID), userJSON, maxAge)

	// Обновляем индекс приватности
	pipe.SAdd(fmt.Sprintf("privacy:%s", user.Privacy), user.ID)

	// Устанавливаем время жизни для ключа приватности
	if maxAge > 0 {
		pipe.Expire(fmt.Sprintf("privacy:%s", user.Privacy), maxAge)
	}

	_, err = pipe.Exec()
	return err
}

// DeleteUser удаляет всю информацию о пользователе из Redis
func (c *GeoCache) DeleteUser(userID int64) error {
	pipe := c.rdb.Pipeline()

	// Получаем информацию о пользователе
	userKey := fmt.Sprintf("user_location:%d", userID)
	userSON, err := c.rdb.Get(userKey).Bytes()
	if err != nil && err != redis.Nil {
		return err
	}

	if err != redis.Nil {
		var user models.User
		err = json.Unmarshal(userSON, &user)
		if err != nil {
			return err
		}

		// Удаляем из геоиндекса
		pipe.ZRem(fmt.Sprintf("h3:%s", user.H3Index), userID)

		// Удаляем из индекса приватности
		pipe.SRem(fmt.Sprintf("privacy:%s", user.Privacy), userID)
	}

	// Удаляем информацию о местоположении пользователя
	pipe.Del(userKey)

	_, err = pipe.Exec()
	return err
}

func (c *GeoCache) GetAllUsersInH3Cell(h3Index string) ([]int64, error) {
	key := fmt.Sprintf("h3:%s", h3Index)

	// Используем GeoRadius с нулевыми координатами и большим радиусом,
	// чтобы получить всех пользователей в этой ячейке H3
	res, err := c.rdb.GeoRadius(key, 0, 0, &redis.GeoRadiusQuery{
		Radius:      1000000, // Большой радиус в км, чтобы охватить всю ячейку
		Unit:        "km",
		WithCoord:   false,
		WithDist:    false,
		WithGeoHash: false,
		Count:       0, // 0 означает "без ограничений"
		Sort:        "ASC",
	}).Result()

	if err != nil {
		if err == redis.Nil {
			// Если ключ не найден, возвращаем пустой слайс
			return []int64{}, nil
		}
		return nil, fmt.Errorf("failed to get users for H3 index %s: %w", h3Index, err)
	}

	userIDs := make([]int64, 0, len(res))
	for _, loc := range res {
		id, err := strconv.ParseInt(loc.Name, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse user ID: %w", err)
		}
		userIDs = append(userIDs, id)
	}

	return userIDs, nil
}

// GetUserLocation получает информацию о местоположении пользователя
func (c *GeoCache) GetUserLocation(userID int64) (*models.User, error) {
	data, err := c.rdb.Get(fmt.Sprintf("user_location:%d", userID)).Bytes()
	if err == redis.Nil {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	var user models.User
	err = json.Unmarshal(data, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// SetPrivacy устанавливает настройки приватности пользователя
func (c *GeoCache) SetPrivacy(ctx context.Context, userID int64, privacy string, maxAge time.Duration) error {
	pipe := c.rdb.Pipeline()

	// Удаляем пользователя из всех существующих индексов приватности
	pipe.SRem("privacy:public", userID)
	pipe.SRem("privacy:friends", userID)
	pipe.SRem("privacy:private", userID)

	// Добавляем пользователя в новый индекс приватности
	pipe.SAdd(fmt.Sprintf("privacy:%s", privacy), userID)

	// Обновляем информацию о местоположении пользователя
	userJSON, err := c.rdb.Get(fmt.Sprintf("user_location:%d", userID)).Bytes()
	if err != nil && err != redis.Nil {
		return err
	}

	if err != redis.Nil {
		var user models.User
		err = json.Unmarshal(userJSON, &user)
		if err != nil {
			return err
		}

		user.Privacy = privacy
		user.UpdatedAt = time.Now()

		updatedJSON, err := json.Marshal(user)
		if err != nil {
			return err
		}

		pipe.Set(fmt.Sprintf("user_location:%d", userID), updatedJSON, maxAge)
	}

	_, err = pipe.Exec()
	return err
}

// GetUsersWithPrivacy получает список пользователей с определенным уровнем приватности
func (c *GeoCache) GetUsersWithPrivacy(privacy string) ([]int64, error) {
	userIDs, err := c.rdb.SMembers(fmt.Sprintf("privacy:%s", privacy)).Result()
	if err != nil {
		return nil, err
	}

	result := make([]int64, len(userIDs))
	for i, idStr := range userIDs {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return nil, err
		}
		result[i] = id
	}

	return result, nil
}
