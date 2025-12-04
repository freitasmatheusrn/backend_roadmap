DB_URL := postgresql://movie_admin:movie_password_123@localhost:5432/movie_reservation?sslmode=disable
MIGRATIONS_DIR := ./internal/database/sql/migrations
SEEDS_DIR := ./internal/database/sql/seeds

.PHONY: migrate-up migrate-down migrate-status migrate-create

migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_URL) up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_URL) down

migrate-down-all:
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_URL) down-to 0

migrate-status:
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_URL) status

migrate-create:
	@read -p "Enter migration name: " name; \
	goose -s -dir $(MIGRATIONS_DIR) create $$name sql

migrate-reset:
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_URL) reset

