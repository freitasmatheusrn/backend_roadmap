-- +goose Up
CREATE TABLE seats (
    id BIGSERIAL PRIMARY KEY,
    theater_id INTEGER NOT NULL,
    row_label VARCHAR(5) NOT NULL,
    seat_number INTEGER NOT NULL,
    seat_type VARCHAR(20),
    is_active BOOLEAN DEFAULT true,
    
    FOREIGN KEY (theater_id) REFERENCES theaters(id) ON DELETE CASCADE,
    UNIQUE (theater_id, row_label, seat_number)
);

CREATE INDEX idx_seats_theater ON seats(theater_id);

-- +goose Down
DROP TABLE seats;