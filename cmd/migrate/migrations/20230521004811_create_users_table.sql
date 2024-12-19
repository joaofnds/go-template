-- +goose Up
CREATE TABLE users(
  id uuid PRIMARY KEY,
  email varchar NOT NULL,
  created_at timestamp NOT NULL
);

-- +goose Down
DROP TABLE users;
