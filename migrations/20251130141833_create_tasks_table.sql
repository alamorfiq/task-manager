-- +goose Up
CREATE TABLE IF NOT EXISTS tasks (
                                     id SERIAL PRIMARY KEY,
                                     title TEXT NOT NULL,
                                     description TEXT DEFAULT '',
                                     is_done BOOLEAN NOT NULL DEFAULT false,
                                     created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
    );

-- +goose Down
DROP TABLE IF EXISTS tasks;