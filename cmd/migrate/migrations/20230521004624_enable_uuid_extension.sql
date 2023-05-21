-- +goose Up
CREATE EXTENSION "uuid-ossp";

-- +goose Down
DROP EXTENSION "uuid-ossp";
