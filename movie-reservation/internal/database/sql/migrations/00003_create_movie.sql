-- +goose Up
CREATE TABLE movies (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    poster_url VARCHAR(500),
    duration_minutes INTEGER NOT NULL,
    release_date DATE,
    rating VARCHAR(10), 
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_duration CHECK (duration_minutes > 0)
);
CREATE INDEX idx_movies_active ON movies(is_active);
CREATE INDEX idx_movies_title ON movies(title);
-- +goose Down
DROP TABLE movies;