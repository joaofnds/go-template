package user_adapter

import (
	"app/user"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var _ user.Repository = &PostgresRepository{}

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db}
}

func (repository *PostgresRepository) CreateUser(ctx context.Context, newUser user.User) error {
	err := repository.db.
		WithContext(ctx).
		Exec("INSERT INTO users(id, email, created_at) VALUES(?, ?, ?)",
			newUser.ID,
			newUser.Email,
			newUser.CreatedAt,
		)

	return gormErr(err)
}
func (repository *PostgresRepository) FindByID(ctx context.Context, id string) (user.User, error) {
	var userFound user.User
	return userFound, gormErr(repository.db.WithContext(ctx).First(&userFound, "id = ?", id))
}

func (repository *PostgresRepository) FindByEmail(ctx context.Context, email string) (user.User, error) {
	var userFound user.User
	return userFound, gormErr(repository.db.WithContext(ctx).First(&userFound, "email = ?", email))
}

func (repository *PostgresRepository) Delete(ctx context.Context, userToDelete user.User) error {
	return gormErr(repository.db.WithContext(ctx).Exec("DELETE FROM users WHERE email = ?", userToDelete.Email))
}

func (repository *PostgresRepository) DeleteAll(ctx context.Context) error {
	return gormErr(repository.db.WithContext(ctx).Exec("DELETE FROM users"))
}

func (repository *PostgresRepository) All(ctx context.Context) ([]user.User, error) {
	var users []user.User
	return users, gormErr(repository.db.WithContext(ctx).Find(&users))
}

func gormErr(result *gorm.DB) error {
	switch {
	case result.Error == nil:
		return nil
	case errors.Is(result.Error, gorm.ErrRecordNotFound):
		return user.ErrNotFound
	default:
		fmt.Printf("\n\n\n%#v\n\n\n", result.Error)
		return user.ErrRepository
	}
}
