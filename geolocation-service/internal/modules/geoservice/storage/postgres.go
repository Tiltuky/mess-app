package storage

import (
	"context"
	"fmt"
	"progekt/dating-app/geolocation-service/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type UserStorage struct {
	db *sqlx.DB
}

func NewUserLocationStorage(db *sqlx.DB) *UserStorage {
	return &UserStorage{
		db: db,
	}
}

func (s *UserStorage) Close() error {
	return s.db.Close()
}

// AddOrUpdateUserLocation добавляет или обновляет местоположение пользователя
func (s *UserStorage) AddOrUpdateUserLocation(ctx context.Context, location *models.User) error {
	if location.ID <= 0 {
		return fmt.Errorf("invalid user ID: %d", location.ID)
	}
	if location.H3Index == "" || len(location.H3Index) != 15 {
		return fmt.Errorf("invalid H3Index: %s", location.H3Index)
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Добавление или обновление в users_table
	query := `
        INSERT INTO users_table (id, privacy, h3_index, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (id) DO UPDATE
        SET privacy = $2, h3_index = $3, updated_at = $5
    `
	_, err = tx.ExecContext(ctx, query, location.ID, location.Privacy, location.H3Index, location.CreatedAt, location.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to add or update user location: %w", err)
	}

	// Добавление записи в location_history
	historyQuery := `
        INSERT INTO location_history (user_id, h3_index, timestamp)
        VALUES ($1, $2, $3)
    `
	_, err = tx.ExecContext(ctx, historyQuery, location.ID, location.H3Index, location.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to add location history: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetUserLocation получает местоположение пользователя по ID
func (s *UserStorage) GetUser(ctx context.Context, userID int64) (*models.User, error) {
	var location models.User
	query := `SELECT * FROM users_table WHERE id = $1`
	err := s.db.GetContext(ctx, &location, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user location: %w", err)
	}
	return &location, nil
}

// GetLocationHistory получает историю местоположений пользователя
func (s *UserStorage) GetLocationHistory(ctx context.Context, userID int64, limit int) (*models.User, error) {
	// Сначала получаем основную информацию о пользователе
	var user models.User
	userQuery := `
        SELECT id, privacy, h3_index, created_at, updated_at
        FROM users_table
        WHERE id = $1
    `
	err := s.db.GetContext(ctx, &user, userQuery, userID)
	if err != nil {
		return &models.User{}, fmt.Errorf("failed to get user info: %w", err)
	}

	// Затем получаем историю местоположений
	historyQuery := `
        SELECT  h3_index, timestamp
        FROM location_history
        WHERE user_id = $1
        ORDER BY timestamp DESC
        LIMIT $2
    `
	err = s.db.SelectContext(ctx, &user.LHistory, historyQuery, userID, limit)
	if err != nil {
		return &models.User{}, fmt.Errorf("failed to get location history: %w", err)
	}

	return &user, nil
}

// ShareLocation начинает шеринг местоположения
func (s *UserStorage) ShareLocation(ctx context.Context, userID int64, sharing *models.LocationSharing) error {
	query := `
        INSERT INTO location_sharing (sharer_id, receiver_id, start_time, end_time)
        VALUES ($1, $2, $3, $4)
    `
	var endTime *time.Time
	if !sharing.EndTime.IsZero() {
		endTime = &sharing.EndTime
	}
	_, err := s.db.ExecContext(ctx, query, userID, sharing.ReceiverID, sharing.StartTime, endTime)
	if err != nil {
		return fmt.Errorf("failed to share location: %w", err)
	}
	return nil
}

// StopSharingLocation прекращает шеринг местоположения
func (s *UserStorage) StopSharingLocation(ctx context.Context, sharerID, receiverID int64) error {
	query := `
        UPDATE location_sharing
        SET end_time = $1
        WHERE sharer_id = $2 AND receiver_id = $3 AND end_time IS NULL
    `
	_, err := s.db.ExecContext(ctx, query, time.Now(), sharerID, receiverID)
	if err != nil {
		return fmt.Errorf("failed to stop sharing location: %w", err)
	}
	return nil
}

// GetActiveSharings получает активные шеринги для пользователя
func (s *UserStorage) GetActiveSharings(ctx context.Context, userID int64) ([]models.LocationSharing, error) {
	query := `
        SELECT id, sharer_id, receiver_id, start_time, COALESCE(end_time, '0001-01-01 00:00:00+00') as end_time FROM location_sharing
        WHERE (sharer_id = $1 OR receiver_id = $1) AND end_time IS NULL
    `
	var sharings []models.LocationSharing
	err := s.db.SelectContext(ctx, &sharings, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active sharings: %w", err)
	}
	return sharings, nil
}

// checkUserIsInActiveSharings проверяет, находится ли пользователь в активных шерингах
func (s *UserStorage) CheckUserIsInActiveSharings(ctx context.Context, responder, target int64) (bool, error) {
	query := `
		SELECT COUNT(*) FROM location_sharing
		WHERE (sharer_id = $1 AND receiver_id = $2) AND end_time IS NULL
		`
	var count int
	err := s.db.GetContext(ctx, &count, query, target, responder)
	if err != nil {
		return false, fmt.Errorf("failed to check user is in active sharings: %w", err)
	}
	return count > 0, nil
}

// UpdatePrivacy обновляет настройки приватности пользователя
func (s *UserStorage) UpdatePrivacy(ctx context.Context, userID int64, privacy string) error {
	query := `
        UPDATE users_table
        SET privacy = $1, updated_at = $2
        WHERE id = $3
    `
	_, err := s.db.ExecContext(ctx, query, privacy, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to update privacy: %w", err)
	}
	return nil
}

// DeleteUser удаляет пользователя и все связанные с ним данные
func (s *UserStorage) DeleteUser(ctx context.Context, userID int64) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Удаление записей из location_sharing
	_, err = tx.ExecContext(ctx, `
        DELETE FROM location_sharing
        WHERE sharer_id = $1 OR receiver_id = $1
    `, userID)
	if err != nil {
		return fmt.Errorf("failed to delete location sharing records: %w", err)
	}

	// Удаление записей из location_history
	_, err = tx.ExecContext(ctx, `
        DELETE FROM location_history
        WHERE user_id = $1
    `, userID)
	if err != nil {
		return fmt.Errorf("failed to delete location history records: %w", err)
	}

	// Удаление пользователя из users_table
	result, err := tx.ExecContext(ctx, `
        DELETE FROM users_table
        WHERE id = $1
    `, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// Проверяем, был ли удален пользователь
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", userID)
	}

	// Фиксируем транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// DeleteHistory удаляет историю местоположений пользователя.
func (s *UserStorage) DeleteHistory(ctx context.Context, userID int64) error {
	query := `
	DELETE FROM location_history
	WHERE user_id = $1
    `
	_, err := s.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete location history: %w", err)
	}
	return nil

}
