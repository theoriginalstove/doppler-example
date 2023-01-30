package app

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"gitlab.com/avocagrow/doppler-example/config"
	"gitlab.com/avocagrow/doppler-example/db"
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
	r.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("at index home page"))
	}))

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
