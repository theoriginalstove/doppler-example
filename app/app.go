package app

import (
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"gitlab.com/steven.t/doppler-example/config"
	"gitlab.com/steven.t/doppler-example/db"
)

type App struct {
	name   string
	db     *db.Instance
	conf   *config.Config
	Server Server
	u      *websocket.Upgrader

	dbMessage chan message
}

func Configure(prefix string, options ...ServerOptionFunc) *App {
	app := &App{
		name: prefix,
		conf: config.Configure(prefix),
		u: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/", home)
	r.HandleFunc("/ws", app.dbInUse)

	r.Post("/_config", app.reloadConfig)

	options = append(options, withHandler(r))
	srv, err := newServer("", options...)
	if err != nil {
		log.Fatal(err)
	}

	app.Server = *srv

	return app
}

func (a *App) reloadConfig(w http.ResponseWriter, r *http.Request) {
	log.Printf("received webhook from doppler to reload config\n")
	a.conf.GetDopplerSecrets()

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

type message struct {
	id   string
	data []byte
}

func (a *App) dbInUse(w http.ResponseWriter, r *http.Request) {
	c, err := a.u.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer c.Close()
	for {

		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
