package app

import (
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"gitlab.com/steven.t/doppler-example/config"
	"gitlab.com/steven.t/doppler-example/db"
)

type App struct {
	name   string
	db     *db.Instance
	conf   *config.Config
	Router http.Handler
}

func Configure(prefix string) *App {
	app := &App{name: prefix}
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/", home)

	r.Post("/_config", reloadConfig)
	app.Router = r
	return app
}

func reloadConfig(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Printf("received webhook from doppler to reload config with context: %v\n", ctx)
	_, err := w.Write([]byte("ok"))
	if err != nil {
		log.Println("err writing to http response")
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		"app/templates/base.html",
		"app/templates/home.html",
	)
	if err != nil {
		log.Println("ERROR: unable to parse home template")
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Printf("ERROR: unable to execute the template: %v\n\n", err)
	}
}
