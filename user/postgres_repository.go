package user

import (
	"context"
	"database/sql"
)

type PostgresRepository struct {
	db Querier
}

type Querier interface {
	QueryRowContext(context.Context, string, ...any) *sql.Row
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db}
}

func (repo *PostgresRepository) CreateUser(ctx context.Context, user User) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO users(name) VALUES($1)", user.Name)
	return err
}

func (repo *PostgresRepository) FindByName(ctx context.Context, name string) (User, error) {
	var user User

	row := repo.db.QueryRowContext(ctx, "SELECT name FROM users WHERE name = $1", name)
	if row.Err() != nil {
		return user, row.Err()
	}

	err := row.Scan(&user.Name)

	return user, err
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
	rows, err := repo.db.QueryContext(ctx, "SELECT name FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return users, err
		}
		users = append(users, User{name})
	}
	return users, nil
}
