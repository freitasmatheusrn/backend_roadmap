package webserver

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/freitasmatheusrn/url-shortening-app/internal/webserver/handlers"
	"github.com/freitasmatheusrn/url-shortening-app/internal/webserver/views"
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
	fileServer := http.FileServer(http.FS(Files))
	r.Handle("/assets/*", fileServer)
	r.Get("/", h.Url.Index)
	r.Get("/form", templ.Handler(views.Form()).ServeHTTP)
	r.Get("/editForm/{id}", h.Url.Edit)
	r.Post("/delete_url/{id}", h.Url.Delete)
	r.Post("/updateUrl/{id}", h.Url.Update)
	r.Post("/createUrl", h.Url.CreateHandler)
	r.Get("/{shortUrl}", h.Url.FowardRequest)

	return r
}
