package storage

import (
	"context"
	"fmt"
	"geolocation-service/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type PayStorage struct {
	db *sqlx.DB
}

func NewPayStorage(db *sqlx.DB) *PayStorage {
	return &PayStorage{
		db: db,
	}
}

func (s *PayStorage) AddCustomer(ctx context.Context, user *models.User) error {
	if user.ID <= 0 {
		return fmt.Errorf("invalid user ID: %d", user.ID)
	}
	if user.CustomerID == "" || user.Name == "" || user.Email == "" {
		return fmt.Errorf("invalid customer details: %v", user)
	}
	query := `
        INSERT INTO customers (user_id, customer_id, name, email,  created_at)
        VALUES ($1, $2, $3, $4, $5,  )
        RETURNING id
    `
	_, err := s.db.ExecContext(ctx, query, user.ID, user.CustomerID, user.Name, user.Email, time.Now())
	if err != nil {
		return fmt.Errorf("failed to add customer: %w", err)
	}
	return nil
}

func (s *PayStorage) UpdateCustomer(ctx context.Context, user *models.User) error {
	if user.ID <= 0 {
		return fmt.Errorf("invalid user ID: %d", user.ID)
	}
	if user.CustomerID == "" || user.Name == "" || user.Email == "" {
		return fmt.Errorf("invalid customer details: %v", user)
	}
	query := `
        UPDATE customers
        SET customer_id = $2, name = $3, email = $4, subscription_end_date = $5
        WHERE user_id = $1
    `
	_, err := s.db.ExecContext(ctx, query, user.ID, user.CustomerID, user.Name, user.Email, user.SubscriptionEndDate)
	if err != nil {
		return fmt.Errorf("failed to update customer: %w", err)
	}
	return nil
}

func (s *PayStorage) GetCustomer(ctx context.Context, userID string) (*models.User, error) {
	query := `
        SELECT * FROM customers WHERE user_id = $1
    `
	var user models.User
	err := s.db.GetContext(ctx, &user, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}
	return &user, nil
}

func (s *PayStorage) DeleteCustomer(ctx context.Context, userID string) error {
	query := `
        DELETE FROM customers WHERE user_id = $1
    `
	_, err := s.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}
	return nil
}

func (s *PayStorage) Close() error {
	return s.db.Close()
}
