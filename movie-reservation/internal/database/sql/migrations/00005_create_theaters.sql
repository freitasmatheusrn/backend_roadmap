-- +goose Up
CREATE TABLE theaters (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    total_seats INTEGER NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT chk_total_seats CHECK (total_seats > 0)
);

-- +goose Down
DROP TABLE theaters;