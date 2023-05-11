package user

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type PostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db}
}

func (repo *PostgresRepository) CreateUser(ctx context.Context, user User) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO users(name) VALUES($1)", user.Name)
	return err
}

func (repo *PostgresRepository) FindByName(ctx context.Context, name string) (User, error) {
	var user User
	return user, repo.db.QueryRowxContext(ctx, "SELECT name FROM users WHERE name = $1", name).StructScan(&user)
}

func (repo *PostgresRepository) Delete(ctx context.Context, user User) error {
	_, err := repo.db.ExecContext(ctx, "DELETE FROM users WHERE name = $1", user.Name)
	return err
}

func (repo *PostgresRepository) DeleteAll(ctx context.Context) error {
	_, err := repo.db.ExecContext(ctx, "DELETE FROM users")
	return err
}

func (repo *PostgresRepository) All(ctx context.Context) ([]User, error) {
	var users []User
	return users, repo.db.SelectContext(ctx, &users, "SELECT name FROM users")
}
