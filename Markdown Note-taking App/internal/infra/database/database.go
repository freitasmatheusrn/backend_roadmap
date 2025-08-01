// internal/infra/database/db.go
package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func NewDB() (*sql.DB, error) {
	connStr := "user=root password=root host=localhost port=5432 dbname=note_app_db sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão com banco: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao conectar com banco: %w", err)
	}
	createNotesTable := `
	CREATE TABLE IF NOT EXISTS notes (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		content BYTEA NOT NULL,
		has_errors BOOLEAN NOT NULL DEFAULT false
	);
	`
	if _, err := db.Exec(createNotesTable); err != nil {
		return nil, fmt.Errorf("erro ao criar tabela notes: %w", err)
	}
	createMatchesTable := `
	CREATE TABLE IF NOT EXISTS matches (
		id SERIAL PRIMARY KEY,
		note_id INTEGER REFERENCES notes(id) ON DELETE CASCADE,
		message TEXT NOT NULL,
		sentence TEXT NOT NULL
	);
	`
	if _, err := db.Exec(createMatchesTable); err != nil {
		return nil, fmt.Errorf("erro ao criar tabela matches: %w", err)
	}
	return db, nil
}
