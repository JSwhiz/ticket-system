package services

import (
	"errors"
	"time"
	"github.com/golang-jwt/jwt/v4"
	"ticket-system/internal/models"
	"ticket-system/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{userRepo, jwtSecret}
}

func (s *AuthService) Login(username, password string) (*models.User, string, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, "", errors.New("пользователь не найден")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", errors.New("неверный пароль")
	}

	// Генерация JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.UserID.String(),
		"role_id": user.RoleID.String(),
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, "", err
	}

	return user, tokenString, nil
}