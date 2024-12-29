-- +goose Up
CREATE TABLE users(
  id UUID PRIMARY KEY,
  email VARCHAR NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE users;
