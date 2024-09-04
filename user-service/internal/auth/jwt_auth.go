package auth

import (
	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	jwtKey         = "secretkey"
	token_exp_time = 5
)

type AuthInterface interface {
	CheckPassword(hashPassword, password string) bool
	GenerateToken(email string) (string, int64, error)
	Manager() *jwtauth.JWTAuth
}

type AuthObj struct {
	AuthManager *jwtauth.JWTAuth
}

func NewAuthObj() *AuthObj {
	return &AuthObj{
		AuthManager: GenerateAuthManager(),
	}
}

func GenerateAuthManager() *jwtauth.JWTAuth {
	AuthManager := jwtauth.New("HS256", []byte(jwtKey), nil)
	return AuthManager
}

func (a *AuthObj) CheckPassword(hashPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err == nil
}

// GenerateToken  токен и expiration time
func (a *AuthObj) GenerateToken(email string) (string, int64, error) {
	claims := map[string]interface{}{"email": email}
	t := time.Now().Add(token_exp_time * time.Minute)
	jwtauth.SetExpiry(claims, t)
	_, token, err := a.AuthManager.Encode(claims)
	if err != nil {
		return "", 0, err
	}

	return token, t.UnixNano(), err
}

func (a *AuthObj) Manager() *jwtauth.JWTAuth {
	return a.AuthManager
}
