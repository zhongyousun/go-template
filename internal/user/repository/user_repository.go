package repository

import (
	"database/sql"
	"go-template/internal/user/model"
)

// UserRepository 定義 contract
type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByEmail(email string) (*model.User, error)
	GetUserByID(id int64) (*model.User, error)
	UpdateUser(user *model.User) error
	DeleteUser(id int64) error
}

// userRepository 是實作，持有 *sql.DB
type userRepository struct {
	DB *sql.DB
}

// NewUserRepository 建立新的 UserRepository 實例
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{DB: db}
}

func (r *userRepository) CreateUser(user *model.User) error {
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

func (r *userRepository) GetUserByID(id int64) (*model.User, error) {
	var user model.User
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

func (r *userRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
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

func (r *userRepository) UpdateUser(user *model.User) error {
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
