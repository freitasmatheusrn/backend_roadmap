package main

import (
	"log"
	"net/http"

	"github.com/freitasmatheusrn/url-shortening-app/internal/database"
	"github.com/freitasmatheusrn/url-shortening-app/internal/repo"
	"github.com/freitasmatheusrn/url-shortening-app/internal/webserver"
	"github.com/freitasmatheusrn/url-shortening-app/internal/webserver/handlers"
)

func main() {
	db, err := database.NewDB("mongodb://root:example@localhost:27017/shortener?authSource=admin", "shortener")

	if err != nil {
		log.Fatal(err)
	}
	urlRepository := repo.NewUrlRepository(db)
	urlHandler := handlers.NewUrlHandler(urlRepository)
	handlers := handlers.NewHandler(urlHandler)
	r := webserver.RegisterRoutes(handlers)
	http.ListenAndServe(":8000", r)
}
