package service

import (
	"context"
	"user-service/internal/models"

	"github.com/brianvoe/gofakeit/v6"
)

func (u *UserServiceObj) CreateFakeUsers(ctx context.Context, countUsers int) error {
	var user models.User
	for i := 0; i < countUsers; i++ {
		gofakeit.Struct(&user)
		_, err := u.CreateUser(ctx, user)
		if err != nil {
			return err
		}
	}
	return nil
}
