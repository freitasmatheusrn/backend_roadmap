package main

import (
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/freitasmatheusrn/movie-reservation/config"
	"github.com/freitasmatheusrn/movie-reservation/internal/auth"
	"github.com/freitasmatheusrn/movie-reservation/internal/genre"
	"github.com/freitasmatheusrn/movie-reservation/internal/movie"
	"github.com/freitasmatheusrn/movie-reservation/internal/reservations"
	"github.com/freitasmatheusrn/movie-reservation/internal/showtimes"
	"github.com/freitasmatheusrn/movie-reservation/internal/user"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/jackc/pgx/v5"
)

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID) // important for rate limiting
	r.Use(middleware.RealIP)    // import for rate limiting and analytics and tracing
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) // recover from crashes

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all good"))
	})
	//users
	tokenGen := auth.NewTokenGenerator(app.config.TokenAuth, int(app.config.JWTExpiresIn))
	userService := user.NewService(user.NewRepo(app.db))
	userHandler := user.NewHandler(userService, tokenGen)
	r.Post("/users", userHandler.Signup)
	r.Post("/sign_in", userHandler.Login)
	movieService := movie.NewService(movie.NewRepo(app.db))
	movieHandler := movie.NewHandler(&movieService)
	r.Get("/movies", movieHandler.List)
	r.Get("/movies/{id}", movieHandler.Show)
	r.Get("/movies_by_genre", movieHandler.ListByGenre)
	r.Get("/movies_search", movieHandler.ListByName)

	// authenticated routes
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(app.config.TokenAuth))
		r.Use(jwtauth.Authenticator)

		genreService := genre.NewService(genre.NewRepo(app.db))
		genreHandler := genre.NewHandler(&genreService)
		showtimeService := showtimes.NewService(showtimes.NewRepo(app.db))
		showtimeHandler := showtimes.NewHandler(&showtimeService)
		reservationService := reservations.NewService(reservations.NewRepo(app.db))
		reservationHandler := reservations.NewHandler(&reservationService)
		r.Post("/reservations/{showtime_id}", reservationHandler.Create)
		r.Get("/reservations/payment_page/{id}", reservationHandler.PaymentPage)
		r.Get("/user_reservations", reservationHandler.List)
		r.Post("/confirm_reservation/{id}", reservationHandler.ConfirmReservation)
		//admin routes
		r.Group(func(r chi.Router) {
			r.Use(auth.AdminMiddleware)
			//genre routes
			r.Post("/genres", genreHandler.Create)
			r.Get("/genres", genreHandler.List)
			//movie routes
			r.Post("/movies", movieHandler.Create)
			//showtimes routes
			r.Post("/showtimes", showtimeHandler.Create)
			r.Get("/showtimes", showtimeHandler.List)
			//reservations routes
			
		})
	})
	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.WebServerPort,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server has started at addr %s", app.config.WebServerPort)

	return srv.ListenAndServeTLS("server.crt", "server.key")
}

type application struct {
	config config.Config
	logger *slog.Logger
	db     *pgx.Conn
}
