package repository

import (
	"context"
	"fmt"
	"os"
	"user-service/internal/models"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserRepoInterface interface {
	SearchUser(ctx context.Context, username string) ([]models.User, error)
	ListUsers(ctx context.Context) ([]models.User, error)
	GetUser(ctx context.Context, id int64) (models.User, error)
	UpdateUser(ctx context.Context, user models.User) error
	DeleteUser(ctx context.Context, id int64) error
	GetUserProfile(ctx context.Context, id int64) (models.UserProfile, error)
	UpdateUserProfile(ctx context.Context, profile models.UserProfile) error
	InsertUser(ctx context.Context, user models.User) (int64, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
}

type UserRepoObj struct {
	db  *gorm.DB
	log *zap.Logger
	es  UserRepoElasticSearchInterface
}

func NewUserRepoObj(log *zap.Logger) *UserRepoObj {
	connString := fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%s sslmode=disable",
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"))
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		log.Error("fail connect to postgres db", zap.String("error", err.Error()))
	}
	es := NewUserRepoElasticSearchObj()
	return &UserRepoObj{
		db:  db,
		log: log,
		es:  es,
	}
}

func (db *UserRepoObj) SearchUser(ctx context.Context, username string) ([]models.User, error) {
	db.log.Info("perform SearchUser query")
	// var user []models.User
	// result := db.db.Where("username ILIKE ?", "%"+username+"%").Find(&user)
	// if result.Error != nil {
	// 	return nil, result.Error
	// }

	user, err := db.es.Search(ctx, username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (db *UserRepoObj) ListUsers(ctx context.Context) ([]models.User, error) {
	db.log.Info("perform ListUsers query")
	var usersList []models.User
	result := db.db.Find(&usersList)
	if result.Error != nil {
		return nil, result.Error
	}
	return usersList, nil
}

func (db *UserRepoObj) GetUser(ctx context.Context, id int64) (models.User, error) {
	db.log.Info("perform GetUser query")
	var user models.User
	result := db.db.First(&user, id)
	if result.Error != nil {
		return models.User{}, result.Error
	}

	return user, nil
}

func (db *UserRepoObj) UpdateUser(ctx context.Context, user models.User) error {
	db.log.Info("perform UpdateUser query")
	result := db.db.Model(&user).Updates(user)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("record not found")
	}

	err := db.es.Update(ctx, &user)
	if err != nil {
		return err
	}
	return nil
}

// В данный момент не используется, польностью удаляет запись из postgres
func (db *UserRepoObj) DeleteUser(ctx context.Context, id int64) error {
	db.log.Info("perform DeleteUser query")
	result := db.db.Delete(models.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("record not found")
	}

	err := db.es.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (db *UserRepoObj) GetUserProfile(ctx context.Context, id int64) (models.UserProfile, error) {
	db.log.Info("perform GetUserProfile query")
	var user models.User
	var profile models.UserProfile
	result := db.db.First(&user, id)
	if result.Error != nil {
		return models.UserProfile{}, result.Error
	}

	profile.Id = user.Id
	profile.Username = user.Username
	profile.FirstName = user.FirstName
	profile.LastName = user.LastName
	profile.Email = user.Email
	profile.Phone = user.Phone
	profile.City = user.City
	profile.Role = user.Role
	profile.AvatarURL = user.AvatarURL

	return profile, nil
}

func (db *UserRepoObj) UpdateUserProfile(ctx context.Context, profile models.UserProfile) error {
	db.log.Info("perform UpdateUserProfile")
	var user models.User
	user.Id = profile.Id
	result := db.db.Model(models.User{}).
		Where("id = ?", profile.Id).
		Select("*").Omit("password", "createdAt", "deletedAt").
		Updates(profile)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("record not found")
	}

	return nil
}

func (db *UserRepoObj) InsertUser(ctx context.Context, user models.User) (int64, error) {
	db.log.Info("perform InsertUser query")
	result := db.db.Create(&user)
	if result.Error != nil {
		return -1, result.Error
	}

	err := db.es.Insert(ctx, &user)
	if err != nil {
		return -1, err
	}
	return user.Id, nil
}

func (db *UserRepoObj) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	db.log.Info("perform GetUserByEmail query")
	var user models.User
	result := db.db.Where("email", email).Find(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}
	if result.RowsAffected == 0 {
		return models.User{}, fmt.Errorf("record not found")
	}
	return user, nil
}
