package user

import (
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

func (repo *PostgresRepository) CreateUser(ctx context.Context, user User) error {
	return gormErr(repo.db.WithContext(ctx).Exec("INSERT INTO users(name) VALUES(?)", user.Name))
}

func (repo *PostgresRepository) FindByName(ctx context.Context, name string) (User, error) {
	var user User
	return user, gormErr(repo.db.WithContext(ctx).First(&user, "name = ?", name))
}

func (repo *PostgresRepository) Delete(ctx context.Context, user User) error {
	return gormErr(repo.db.WithContext(ctx).Exec("DELETE FROM users WHERE name = ?", user.Name))
}

func (repo *PostgresRepository) DeleteAll(ctx context.Context) error {
	return gormErr(repo.db.WithContext(ctx).Exec("DELETE FROM users"))
}

func (repo *PostgresRepository) All(ctx context.Context) ([]User, error) {
	var users []User
	return users, gormErr(repo.db.WithContext(ctx).Find(&users))
}

func gormErr(result *gorm.DB) error {
	switch {
	case result.Error == nil:
		return nil
	case errors.Is(result.Error, gorm.ErrRecordNotFound):
		return ErrNotFound
	default:
		return ErrRepository
	}
}
