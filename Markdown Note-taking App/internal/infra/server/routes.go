package server

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/freitasmatheusrn/markdown-note-taking-app/internal/infra/server/handlers"
	"github.com/freitasmatheusrn/markdown-note-taking-app/internal/infra/web"
	"github.com/freitasmatheusrn/markdown-note-taking-app/internal/infra/web/views"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func RegisterRoutes(h *handlers.Handlers) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	fileServer := http.FileServer(http.FS(web.Files))
	r.Handle("/assets/*", fileServer)
	r.Get("/", h.Note.Home)
	r.Get("/form", templ.Handler(views.UploadNote()).ServeHTTP)
	r.Get("/notes/{id}", h.Note.ShowNote)
	r.Post("/notes/upload", h.Note.UploadNote)

	return r
}
