package rhandlers

import (
	"fmt"
	"myapp2/pkg/config"
	"myapp2/pkg/handlers"
	"myapp2/pkg/helpers"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig
var session *scs.SessionManager

func SetConfig(config *config.AppConfig) {

	app = config

}

func Auth(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsAuthenticated(r) {
			app.Session.Put(r.Context(), "error", "log in first")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func WriteToConsole(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("hit the page")
		next.ServeHTTP(w, r)
	})

}

// NoSurf adds Csrf protection to all post requests
func NoSurf(next http.Handler) http.Handler {

	csrfHandler := nosurf.New(next)
	//fmt.Println(csrfHandler)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})
	//fmt.Println(csrfHandler)
	return csrfHandler

}

//SessionLoad loads and saves session on every request
func SessionLoad(next http.Handler) http.Handler {

	return app.Session.LoadAndSave(next)

}

func Routes2(app *config.AppConfig) http.Handler {

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	//mux.Use(WriteToConsole)
	mux.Use(NoSurf)

	mux.Use(SessionLoad)

	mux.Get("/", handlers.Home)
	mux.Get("/about", handlers.About)
	mux.Get("/index", handlers.Index)
	mux.Get("/rooms/room1", handlers.Rooms)
	mux.Get("/rooms/room2", handlers.Rooms)
	mux.Get("/book", handlers.Reservation)
	mux.Post("/book", handlers.PostReservation)

	mux.Get("/reservation", handlers.MakeReservations)

	mux.Post("/reservation", handlers.PostMakeReservations)

	mux.Get("/choose_room/{id}", handlers.ChooseRoom)

	mux.Post("/book-json", handlers.AvailibilityJson)

	mux.Get("/reservation-sum", handlers.ReservationSummary)

	mux.Get("/user/login", handlers.ShowLogin)

	mux.Post("/user/login", handlers.PostShowLogin)

	mux.Get("/user/logout", handlers.Logout)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Route("/admin", func(r chi.Router) {
		//r.Use(Auth)
		r.Get("/dashboard", handlers.AdminDashboard)
		r.Get("/reservations-new", handlers.AdminNewReservations)
		r.Get("/reservations-all", handlers.AdminAllReservations)

		r.Get("/reservations-calendar", handlers.AdminReservationsCalander)
		r.Post("/reservations-calendar", handlers.AdminPostReservationsCalander)

		r.Get("/process-reservation/{src}/{id}", handlers.AdminProcessReservation)
		r.Get("/delete-reservation/{src}/{id}", handlers.AdminDeleteReservation)

		r.Get("/reservations/{src}/{id}/show", handlers.AdminShowReservation)
		r.Post("/reservations/{src}/{id}/show", handlers.AdminPostShowReservation)

	})

	return mux
}
