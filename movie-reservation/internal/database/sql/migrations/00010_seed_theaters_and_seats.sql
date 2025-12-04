-- +goose Up
-- Inserir salas
INSERT INTO theaters (name, total_seats)
VALUES 
  ('Sala 1', 100),
  ('Sala 2', 80),
  ('Sala 3', 120);

-- Criar assentos automaticamente
-- Sala 1 → 10 fileiras A-J, 10 assentos cada
WITH theater AS (
    SELECT id FROM theaters WHERE name = 'Sala 1'
)
INSERT INTO seats (theater_id, row_label, seat_number, seat_type)
SELECT 
    t.id,
    chr(64 + r.row) AS row_label,
    s.seat AS seat_number,
    'standard'
FROM theater t
CROSS JOIN generate_series(1, 10) AS r(row)
CROSS JOIN generate_series(1, 10) AS s(seat);

-- Sala 2 → 8 fileiras A-H, 10 assentos cada
WITH theater AS (
    SELECT id FROM theaters WHERE name = 'Sala 2'
)
INSERT INTO seats (theater_id, row_label, seat_number, seat_type)
SELECT 
    t.id,
    chr(64 + r.row) AS row_label,
    s.seat AS seat_number,
    'standard'
FROM theater t
CROSS JOIN generate_series(1, 8) AS r(row)
CROSS JOIN generate_series(1, 10) AS s(seat);

-- Sala 3 → 12 fileiras A-L, 10 assentos cada
WITH theater AS (
    SELECT id FROM theaters WHERE name = 'Sala 3'
)
INSERT INTO seats (theater_id, row_label, seat_number, seat_type)
SELECT 
    t.id,
    chr(64 + r.row) AS row_label,
    s.seat AS seat_number,
    CASE 
        WHEN r.row <= 2 THEN 'vip'
        ELSE 'standard'
    END AS seat_type
FROM theater t
CROSS JOIN generate_series(1, 12) AS r(row)
CROSS JOIN generate_series(1, 10) AS s(seat);


-- +goose Down
-- Remover assentos e salas
DELETE FROM seats WHERE theater_id IN (SELECT id FROM theaters WHERE name IN ('Sala 1','Sala 2','Sala 3'));
DELETE FROM theaters WHERE name IN ('Sala 1','Sala 2','Sala 3');
