package service

import (
	"fmt"
	"os"
	"time"

	"go-template/internal/db"
	orderModel "go-template/internal/order/model"
	orderrepo "go-template/internal/order/repository"
	userModel "go-template/internal/user/model"
	userrepo "go-template/internal/user/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// UserService is responsible for user-related operations
type UserService struct {
	Repo      userrepo.UserRepository
	OrderRepo orderrepo.OrderRepository
	txManager db.TransactionManager
}

func NewUserService(repo userrepo.UserRepository, orderRepo orderrepo.OrderRepository) *UserService {
	return &UserService{Repo: repo, OrderRepo: orderRepo}
}

func NewUserServiceWithTx(repo userrepo.UserRepository, orderRepo orderrepo.OrderRepository, txManager db.TransactionManager) *UserService {
	return &UserService{
		Repo:      repo,
		OrderRepo: orderRepo,
		txManager: txManager,
	}
}

func (s *UserService) GetUserByID(id int64) (*userModel.User, error) {
	return s.Repo.GetUserByID(id)
}

func (s *UserService) GetUserByIDWithCache(id int64) (*userModel.User, error) {
	return s.Repo.GetUserByIDWithCache(id)
}

// GetUserWithOrders returns user and their orders by userID
func (s *UserService) GetUserWithOrders(userID int64) (*userModel.UserWithOrders, error) {
	usr, err := s.Repo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("get user failed: %w", err)
	}
	if usr == nil {
		return nil, nil
	}
	orders, err := s.OrderRepo.GetOrdersByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("get orders failed: %w", err)
	}
	return &userModel.UserWithOrders{
		User:   usr,
		Orders: orders,
	}, nil
}

func (s *UserService) RegisterUser(user *userModel.User) error {
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

func (s *UserService) UpdateUser(user *userModel.User) error {
	return s.Repo.UpdateUser(user)
}

func (s *UserService) DeleteUser(id int64) error {
	return s.Repo.DeleteUser(id)
}

func (s *UserService) RegisterUserWithOrder(user *userModel.User, order *orderModel.Order) error {
	uow, err := s.txManager.Begin()
	if err != nil {
		return err
	}
	defer uow.Rollback() // Rollback if any error occurs

	if err := uow.UserRepo().CreateUser(user); err != nil {
		return err
	}

	order.UserID = int64(user.ID)
	if err := uow.OrderRepo().CreateOrder(order); err != nil {
		return err
	}

	return uow.Commit() // Commit only if all succeed
}
