package service

import (
	"context"
	"fmt"
	"os"
	"strings"
	"user-service/internal/auth"
	"user-service/internal/models"

	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	goauth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
)

type GoogleServiceInterface interface {
	Authenticate(ctx context.Context, accessToken, pushToken string) (string, int64, error)
}

type GoogleService struct {
	User   UserServiceInterface
	Auth   auth.AuthInterface
	logger *zap.Logger
}

func NewGoogleService(userService UserServiceInterface, logger *zap.Logger) GoogleServiceInterface {
	return &GoogleService{User: userService, Auth: auth.NewAuthObj(), logger: logger}
}

func (s *GoogleService) Authenticate(ctx context.Context, accessToken string, pushToken string) (string, int64, error) {
	config := oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/contacts.readonly"},
		Endpoint:     google.Endpoint,
	}

	oauthToken := &oauth2.Token{
		AccessToken: accessToken,
	}

	client := config.Client(ctx, oauthToken)
	_, err := people.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		s.logger.Error("Failed to create google people api service:", zap.Error(err))
	}

	oauth2Service, err := goauth2.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		s.logger.Error("Failed to create google oauth2 api service:", zap.Error(err))
	}

	userInfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		s.logger.Error("Failed to get user info:", zap.Error(err))
	}

	user, err := s.User.GetUserByEmail(ctx, userInfo.Email)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		s.logger.Error("Failed to get user:", zap.Error(err))
	}

	if user.Id == 0 {
		username := fmt.Sprintf("%s_%s", userInfo.GivenName, userInfo.FamilyName)
		newUser := models.User{
			Username:  username,
			FirstName: userInfo.GivenName,
			LastName:  userInfo.FamilyName,
			Email:     userInfo.Email,
			AvatarURL: userInfo.Picture,
		}

		_, err := s.User.CreateUser(ctx, newUser)
		if err != nil {
			s.logger.Error("Failed to register user:", zap.Error(err))
		}

	}

	token, exp, err := s.Auth.GenerateToken(userInfo.Email)

	if err != nil {
		return "", 0, err
	}

	// Send notification event to broker

	return token, exp, nil
}
