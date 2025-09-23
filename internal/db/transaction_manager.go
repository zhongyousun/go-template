package db

import (
	orderRepo "go-template/internal/order/repository"
	userRepo "go-template/internal/user/repository"

	"gorm.io/gorm"
)

// TransactionManager is responsible for starting transactions
type TransactionManager interface {
	Begin() (UnitOfWork, error)
}

// UnitOfWork represents a transaction scope
type UnitOfWork interface {
	UserRepo() userRepo.UserRepository
	OrderRepo() orderRepo.OrderRepository
	Commit() error
	Rollback() error
}

// gormTxManager implements TransactionManager using GORM
type gormTxManager struct {
	db *gorm.DB
}

// NewTransactionManager creates a new TransactionManager
func NewTransactionManager(db *gorm.DB) TransactionManager {
	return &gormTxManager{db: db}
}

// Begin starts a transaction and returns a UnitOfWork
func (m *gormTxManager) Begin() (UnitOfWork, error) {
	tx := m.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &gormUnitOfWork{
		tx:        tx,
		userRepo:  userRepo.NewUserRepository(tx),   // 把 tx 傳進 repository
		orderRepo: orderRepo.NewOrderRepository(tx), // 同上
	}, nil
}

// gormUnitOfWork is a GORM implementation of UnitOfWork
type gormUnitOfWork struct {
	tx        *gorm.DB
	userRepo  userRepo.UserRepository
	orderRepo orderRepo.OrderRepository
}

func (u *gormUnitOfWork) UserRepo() userRepo.UserRepository {
	return u.userRepo
}

func (u *gormUnitOfWork) OrderRepo() orderRepo.OrderRepository {
	return u.orderRepo
}

func (u *gormUnitOfWork) Commit() error {
	return u.tx.Commit().Error
}

func (u *gormUnitOfWork) Rollback() error {
	return u.tx.Rollback().Error
}
