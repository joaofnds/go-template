package adapter

import (
	"app/user"
	"context"
	"errors"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db}
}

func (repository *PostgresRepository) CreateUser(ctx context.Context, newUser user.User) error {
	return gormErr(repository.db.WithContext(ctx).Exec("INSERT INTO users(name) VALUES(?)", newUser.Name))
}

func (repository *PostgresRepository) FindByName(ctx context.Context, name string) (user.User, error) {
	var userFound user.User
	return userFound, gormErr(repository.db.WithContext(ctx).First(&userFound, "name = ?", name))
}

func (repository *PostgresRepository) Delete(ctx context.Context, userToDelete user.User) error {
	return gormErr(repository.db.WithContext(ctx).Exec("DELETE FROM users WHERE name = ?", userToDelete.Name))
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
		return user.ErrRepository
	}
}
