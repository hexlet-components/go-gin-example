-- +goose Up
CREATE TABLE IF NOT EXISTS articles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS articles;
