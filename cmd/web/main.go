package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"myapp2/pkg/config"
	"myapp2/pkg/driver"
	"myapp2/pkg/handlers"
	"myapp2/pkg/helpers"
	"myapp2/pkg/models"
	"myapp2/pkg/renderers"
	"myapp2/pkg/rhandlers"
	"net/http"
	//"net/smtp"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
)

var app config.AppConfig
var infoLog *log.Logger
var errorLog *log.Logger

func main() {

	db, _ := run()
	defer db.SQL.Close()

	defer close(app.MailChan)

	listenForMail()

	msg := models.MailData{
		To:      "m@c.com",
		From:    "abc@c.com",
		Subject: "hi i am here",
		Content: "well done for your report lets play cricket",
	}

	app.MailChan <- msg

	/*

		from := "me@here.com"
		auth := smtp.PlainAuth("",from,"","localhost")
		err := smtp.SendMail("localhost:1025",auth,from,[]string{"timeisgod12@gmail.com"},[]byte("Hello World"))
		if err!=nil{
			log.Println(err)
		}*/

	fmt.Println("started version 2")

	server := &http.Server{
		Addr:    ":8080",
		Handler: rhandlers.Routes2(&app),
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

func run() (*driver.DB, error) {

	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	// read flags

	mailChan := make(chan models.MailData)

	app.MailChan = mailChan

	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session := scs.New()

	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to database

	log.Println("connecting to database")

	db, err := driver.ConnectSql("host=localhost port=5432 dbname=bookings user=cm1100 password=srsscthr")
	if err != nil {
		log.Fatal(err)
	}

	handlers.SetConfig(&app)
	handlers.SetRepository(db)
	rhandlers.SetConfig(&app)
	renderers.SetAppConfig(&app)
	helpers.SetConfig(&app)

	return db, nil

}
