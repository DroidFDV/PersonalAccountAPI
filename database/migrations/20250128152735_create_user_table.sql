-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
				id SERIAL PRIMARY KEY,
				login VARCHAR(100) NOT NULL,
				password VARCHAR(255) NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd