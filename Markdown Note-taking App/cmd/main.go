package main

import (
	"net/http"

	"github.com/freitasmatheusrn/markdown-note-taking-app/internal/infra/database"
	"github.com/freitasmatheusrn/markdown-note-taking-app/internal/infra/server"
	"github.com/freitasmatheusrn/markdown-note-taking-app/internal/infra/server/handlers"
	"github.com/freitasmatheusrn/markdown-note-taking-app/internal/repositories"
)

func main() {
	db, err := database.NewDB()
	if err != nil {
		panic(err)
	}
	noteRepository := repositories.NewNoteRepository(db)
	noteHandler := handlers.NewNoteHandler(noteRepository)
	handlersHandlers := handlers.NewHandlers(noteHandler)
	r := server.RegisterRoutes(handlersHandlers)
	http.ListenAndServe(":8000", r)
}
