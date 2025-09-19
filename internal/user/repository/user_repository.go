package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	user "go-template/internal/user/model"
	"go-template/pkg/redisclient"
	"time"
)

// UserRepository defines the contract for user data access
type UserRepository interface {
	CreateUser(user *user.User) error
	GetUserByEmail(email string) (*user.User, error)
	GetUserByID(id int64) (*user.User, error)
	GetUserByIDWithCache(id int64) (*user.User, error)
	UpdateUser(user *user.User) error
	DeleteUser(id int64) error
}

// userRepository is the implementation, holding *sql.DB
type userRepository struct {
	DB *sql.DB
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{DB: db}
}

func (r *userRepository) CreateUser(user *user.User) error {
	err := r.DB.QueryRow(
		`INSERT INTO "user" ("name", "email", "password", "role")
		VALUES ($1, $2, $3, $4)
		RETURNING "id", "createdAt"`,
		user.Name, user.Email, user.Password, user.Role,
	).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) GetUserByID(id int64) (*user.User, error) {
	var user user.User
	err := r.DB.QueryRow(
		`SELECT "id", "name", "email", "password", "createdAt", "role"
		FROM "user" WHERE "id" = $1`,
		id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.Role)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByIDWithCache(id int64) (*user.User, error) {
	cacheKey := fmt.Sprintf("user:%d", id)

	// 1. Try Redis cache first
	if redisclient.Rdb != nil {
		if val, err := redisclient.Rdb.Get(context.Background(), cacheKey).Result(); err == nil {
			var u user.User
			if err := json.Unmarshal([]byte(val), &u); err == nil {
				return &u, nil // cache hit
			}
		}
	}

	// 2. Query DB if cache miss
	var u user.User
	err := r.DB.QueryRow(
		`SELECT "id", "name", "email", "password", "createdAt", "role" 
		FROM "user" WHERE "id" = $1`,
		id,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.CreatedAt, &u.Role)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	// 3. Save result into Redis (cache for 10 minutes)
	if redisclient.Rdb != nil {
		if bytes, err := json.Marshal(u); err == nil {
			_ = redisclient.Rdb.Set(context.Background(), cacheKey, bytes, 10*time.Minute).Err()
		}
	}

	return &u, nil
}

func (r *userRepository) GetUserByEmail(email string) (*user.User, error) {
	var user user.User
	err := r.DB.QueryRow(
		`SELECT "id", "name", "email", "password", "role"
          FROM "user"
          WHERE "email" = $1`,
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateUser(user *user.User) error {
	res, err := r.DB.Exec(
		`UPDATE "user" SET "name"=$1, "email"=$2, "password"=$3, "role"=$4 WHERE "id"=$5`,
		user.Name, user.Email, user.Password, user.Role, user.ID,
	)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *userRepository) DeleteUser(id int64) error {
	res, err := r.DB.Exec(`DELETE FROM "user" WHERE "id"=$1`, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
