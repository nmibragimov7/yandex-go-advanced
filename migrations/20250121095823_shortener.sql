-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS shortener (
                                         id SERIAL PRIMARY KEY,
                                         short_url VARCHAR(10) NOT NULL,
    original_url VARCHAR(100) UNIQUE NOT NULL
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE shortener;
-- +goose StatementEnd