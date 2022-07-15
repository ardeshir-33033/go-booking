package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/ardeshir-33033/go-booking/pkg/config"
	"github.com/ardeshir-33033/go-booking/pkg/handlers/handlers"
	"github.com/ardeshir-33033/go-booking/pkg/render"
	"log"
	"net/http"
	"time"
)

const portNumber = ":8082"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot CreateTemplateCache")
	}
	app.TemplateCache = tc
	app.UseCache = false
	render.NewTemplates(&app)

	repo := handlers.NewRepo(&app)
	handlers.NewHandler(repo)

	//http.HandleFunc("/", handlers.Repo.Home)
	//http.HandleFunc("/about", handlers.Repo.About)
	//http.HandleFunc("/divide", handlers.Divide)

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))

	//_ = http.ListenAndServe(portNumber, nil)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)

}
