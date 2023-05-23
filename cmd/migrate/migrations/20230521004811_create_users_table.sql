-- +goose Up
CREATE TABLE users(
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  name varchar NOT NULL
);

-- +goose Down
DROP TABLE users;
