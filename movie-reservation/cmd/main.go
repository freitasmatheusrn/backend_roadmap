package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/freitasmatheusrn/movie-reservation/config"
	"github.com/freitasmatheusrn/movie-reservation/internal/database"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	dsn := fmt.Sprintf(
		"%s://%s:%s@localhost:5432/%s", config.DBDriver, config.DBUser, config.DBPassword, config.DBName)
	db, err := database.InitPostgres(dsn)
	if err != nil {
		panic("error starting db")
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	app := application{
		config: *config,
		logger: logger,
		db:     db,
	}

	if err := app.run(app.mount()); err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}

}
