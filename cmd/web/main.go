package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/SujithSubhash/bookings/internal/config"
	"github.com/SujithSubhash/bookings/internal/drivers"
	"github.com/SujithSubhash/bookings/internal/handlers"
	"github.com/SujithSubhash/bookings/internal/helpers"
	"github.com/SujithSubhash/bookings/internal/models"
	"github.com/SujithSubhash/bookings/internal/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":3333"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)

	fmt.Printf("Starting application on port number%s", portNumber)

	//_ = http.ListenAndServe(portNumber, nil)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (*drivers.DB, error) {

	//What am i going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	//change this to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog
	//setup the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	//connect to database

	log.Println("Connecting to database....")
	db, err := drivers.ConnectSQL("host=localhost port=5432 dbname=bookings user=sujithsubhash password=")
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}
	log.Println("Connected to database!")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
		return nil, err
	}
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)

	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
