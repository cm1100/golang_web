package handlers

import (
	"encoding/gob"
	"myapp2/pkg/config"
	"myapp2/pkg/models"
	"myapp2/pkg/renderers"

	"html/template"
	"log"

	//"myapp2/pkg/renderers"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"

	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
)

//var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../templates"
var functions = template.FuncMap{}

var app2 config.AppConfig

func TestMain(m *testing.M) {
	gob.Register(models.Reservation{})

	app2.InProduction = false

	session = scs.New()

	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app2.InProduction

	app2.Session = session

	mailChan := make(chan models.MailData)
	app2.MailChan = mailChan
	defer close(mailChan)

	listenForMail()

	log.Println("connecting to database")

	//db2,err:= driver.ConnectSql("host=localhost port=5432 dbname=bookings user=cm1100 password=srsscthr")
	//if err!=nil{
	//	log.Fatal(err)
	//}

	//fmt.Println(db,"connected maybe")

	SetConfig(&app2)
	renderers.SetAppConfig(&app2)
	SetTestRepositery()
	renderers.SetPath("./../../templates")

	os.Exit(m.Run())

}

func getRoutes() http.Handler {

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	//mux.Use(WriteToConsole)
	//mux.Use(NoSurf)

	mux.Use(SessionLoad)

	mux.Get("/", Home)
	mux.Get("/about", About)
	mux.Get("/index", Index)
	mux.Get("/rooms/room1", Rooms)
	mux.Get("/rooms/room2", Rooms)
	mux.Get("/book", Reservation)
	mux.Post("/book", PostReservation)

	mux.Post("/book-json", AvailibilityJson)

	mux.Get("/reservation", MakeReservations)

	mux.Post("/reservation", PostMakeReservations)

	mux.Get("/choose_room/{id}", ChooseRoom)

	mux.Get("/reservation-sum", ReservationSummary)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux

}

func NoSurf(next http.Handler) http.Handler {

	csrfHandler := nosurf.New(next)
	//fmt.Println(csrfHandler)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Secure:   app2.InProduction,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})
	//fmt.Println(csrfHandler)
	return csrfHandler

}

//SessionLoad loads and saves session on every request
func SessionLoad(next http.Handler) http.Handler {

	return app2.Session.LoadAndSave(next)

}

func listenForMail() {
	go func() {
		for {
			_ = <-app2.MailChan
		}
	}()
}
