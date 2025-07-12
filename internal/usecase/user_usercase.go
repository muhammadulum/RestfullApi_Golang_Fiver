package usecase

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"learn_restful_api_golang/internal/domain"
	"learn_restful_api_golang/pkg/utils"
)

type userUsecase struct {
	repo domain.UserRepository
}

func NewUserUseCase(r domain.UserRepository) domain.AuthUseCase {
	return &userUsecase{repo: r}
}

func (u *userUsecase) Login(email, password string) (string, string, error) {
	user, err := u.repo.FindByEmail(email)
	if err != nil {
		return "", "", errors.New("user not found")
	}
	if !utils.CheckPasswordHash(password, user.Password) {
		return "", "", errors.New("invalid credentials")
	}

	accessToken, err := u.generateToken(user, time.Minute*15)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := u.generateToken(user, time.Hour*24*7)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (u *userUsecase) generateToken(user *domain.User, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(duration).Unix(),
	})
	secret := os.Getenv("JWT_SECRET")
	return token.SignedString([]byte(secret))
}

func (u *userUsecase) RefreshToken(oldToken string) (string, error) {
	parsed, err := jwt.Parse(oldToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if claims, ok := parsed.Claims.(jwt.MapClaims); ok && parsed.Valid {
		user := &domain.User{
			ID:    uint(claims["user_id"].(float64)),
			Email: claims["email"].(string),
			Role:  claims["role"].(string),
		}
		return u.generateToken(user, time.Minute*15)
	}
	return "", err
}

func (u *userUsecase) Register(name, email, password string) error {
	hashed, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	user := &domain.User{
		Name:     name,
		Email:    email,
		Password: hashed,
		Role:     "user",
	}
	return u.repo.Create(user)
}