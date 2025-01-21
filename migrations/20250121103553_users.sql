-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users
(
    id SERIAL PRIMARY KEY
);
ALTER TABLE shortener
    ADD COLUMN user_id INTEGER REFERENCES users(id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
ALTER TABLE shortener
    DROP COLUMN user_id;
-- +goose StatementEnd
