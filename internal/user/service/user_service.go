package service

import (
	"fmt"
	"os"
	"time"

	"go-template/internal/user/model"
	"go-template/internal/user/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type UserService struct {
	Repo repository.UserRepository // 注意這裡依賴的是 interface
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) RegisterUser(user *model.User) error {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hashedPassword)
	return s.Repo.CreateUser(user)
}

func (s *UserService) LoginUser(email, password string) (string, error) {
	user, err := s.Repo.GetUserByEmail(email)
	if err != nil {
		return "", fmt.Errorf("failed to query user: %w", err)
	}
	if user == nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", fmt.Errorf("invalid email or password")
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, nil
}

func (s *UserService) UpdateUser(user *model.User) error {
	return s.Repo.UpdateUser(user)
}

func (s *UserService) DeleteUser(id int64) error {
	return s.Repo.DeleteUser(id)
}

func (s *UserService) GetUserByID(id int64) (*model.User, error) {
	return s.Repo.GetUserByID(id)
}
