package service

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"user-service/internal/auth"
	"user-service/internal/models"
	"user-service/internal/repository"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const uploads = "uploads/"

type UserServiceInterface interface {
	CreateUser(ctx context.Context, user models.User) (int64, error)
	UpdateUser(ctx context.Context, user models.User) error
	GetUser(ctx context.Context, id int64) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	DeleteUser(ctx context.Context, id int64) error
	GetUserProfile(ctx context.Context, id int64) (models.UserProfile, error)
	UpdateUserProfile(ctx context.Context, profile models.UserProfile) error
	ListUsers(ctx context.Context) ([]models.User, error)
	UploadAvatar(ctx context.Context, img []byte, filename string, id int64) error
	AuthenticateUser(ctx context.Context, email string, password string) (string, int64, error)
	SearchUser(ctx context.Context, username string) ([]models.User, error)
}

type UserServiceObj struct {
	repo repository.UserRepoInterface
	auth auth.AuthInterface
	log  *zap.Logger
}

func NewUserServiceObj(repo repository.UserRepoInterface, log *zap.Logger) *UserServiceObj {
	return &UserServiceObj{
		repo: repo,
		auth: auth.NewAuthObj(),
		log:  log,
	}
}

func (s *UserServiceObj) AuthenticateUser(ctx context.Context, email string, password string) (string, int64, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", 0, err
	}
	if !s.auth.CheckPassword(user.Password, password) {
		s.log.Info("auth failure, wrong password")
		return "", 0, fmt.Errorf("wrong password")
	}

	token, exp, err := s.auth.GenerateToken(email)
	if err != nil {
		return "", 0, err
	}

	return token, exp, nil
}

func (s *UserServiceObj) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s *UserServiceObj) CreateUser(ctx context.Context, user models.User) (int64, error) {
	s.log.Info("hash password")
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return -1, err
	}
	user.Password = string(hashedPass)
	user.CreatedAt = time.Now().UnixNano()
	user.DeletedAt = user.CreatedAt - 1
	id, err := s.repo.InsertUser(ctx, user)

	return id, err
}

func (s *UserServiceObj) UpdateUser(ctx context.Context, user models.User) error {
	err := s.repo.UpdateUser(ctx, user)

	return err
}

func (s *UserServiceObj) GetUser(ctx context.Context, id int64) (models.User, error) {
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// Устанавливает DeletedAt, не удалая саму запись
func (s *UserServiceObj) DeleteUser(ctx context.Context, id int64) error {
	user := models.User{
		Id:        id,
		DeletedAt: time.Now().UnixNano(),
	}
	err := s.repo.UpdateUser(ctx, user)

	return err
}

func (s *UserServiceObj) GetUserProfile(ctx context.Context, id int64) (models.UserProfile, error) {
	var profile models.UserProfile
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return models.UserProfile{}, err
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

	return profile, err

}

func (s *UserServiceObj) UpdateUserProfile(ctx context.Context, profile models.UserProfile) error {
	var user models.User

	user.Id = profile.Id
	user.Username = profile.Username
	user.FirstName = profile.FirstName
	user.LastName = profile.LastName
	user.Email = profile.Email
	user.Phone = profile.Phone
	user.City = profile.City
	user.Role = profile.Role
	user.AvatarURL = profile.AvatarURL

	err := s.repo.UpdateUser(ctx, user)

	return err

}

func (s *UserServiceObj) SearchUser(ctx context.Context, username string) ([]models.User, error) {
	user, err := s.repo.SearchUser(ctx, username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserServiceObj) ListUsers(ctx context.Context) ([]models.User, error) {
	users, err := s.repo.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		err := s.CreateFakeUsers(ctx, 5)
		if err != nil {
			return nil, err
		}
		users, err = s.repo.ListUsers(ctx)
		if err != nil {
			return nil, err
		}
	}

	return users, nil
}

func (s *UserServiceObj) UploadAvatar(ctx context.Context, img []byte, filename string, id int64) error {
	splitted := strings.Split(filename, ".")
	if len(splitted) == 1 {
		return fmt.Errorf("wrong file extension")
	}
	extension := splitted[len(splitted)-1]
	filePath := uploads + strconv.FormatInt(id, 10) + "." + extension

	err := os.WriteFile(filePath, img, 0666)
	if err != nil {
		return fmt.Errorf("uploadAvatar err:%w", err)
	}

	user := models.User{
		Id:        id,
		AvatarURL: filePath,
	}
	fmt.Println(id, filePath)
	err = s.repo.UpdateUser(ctx, user)

	return err
}
