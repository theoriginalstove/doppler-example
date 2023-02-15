package app

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"gitlab.com/steven.t/doppler-example/config"
	"gitlab.com/steven.t/doppler-example/db"
)

var ensureScheme = flag.Bool("schema", true, "If the db should ensure the table schema exists")

type App struct {
	name   string
	Db     *db.Instance
	conf   *config.Config
	Server Server
	u      *websocket.Upgrader

	// connection hub
	hub *Hub
	// data hub
	dataHub *Hub
}

func Configure(prefix string, options ...ServerOptionFunc) *App {
	app := &App{
		name: prefix,
		conf: config.Configure(prefix),
		u:    &websocket.Upgrader{},
	}
	err := app.conf.GetDopplerSecrets()
	if err != nil {
		log.Fatal("unable to get doppler secrets")
	}

	app.Db = db.Configure(true, "", app.conf)

	app.hub = newHub()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/", app.home)
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		app.dbInUse(app.hub, w, r)
	})
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "app/assets/favicon.ico")
	})

	r.Post("/_config", app.reloadConfig)

	options = append(options, withHandler(r))
	srv, err := newServer("", options...)
	if err != nil {
		log.Fatal(err)
	}
	go app.hub.run()

	app.Server = *srv

	return app
}

func (a *App) reloadConfig(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := a.conf.GetDopplerSecrets()
	if err != nil {
		http.Error(w, fmt.Sprintf("%d - Server Error - unable to refresh secrets", http.StatusInternalServerError), http.StatusInternalServerError)
	}

	// Reset the db connection
	connStr := a.conf.Secrets["ROACH_CONN"]
	err = a.Db.SetNewConnection(ctx, connStr)
	if err != nil {
		log.Printf("error setting new db connection: %w", err)
		return
	}

	db := a.conf.Secrets["ROACH_DB"]
	a.hub.broadcast <- []byte(db)

	_, err = w.Write([]byte("ok"))
	if err != nil {
		log.Println("err writing to http response")
	}
}

type home struct {
	ConnectedTo string
	Plants      []db.Plant
}

func (a *App) home(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		"app/templates/base.html",
		"app/templates/home.html",
	)
	if err != nil {
		log.Println("ERROR: unable to parse home template")
	}
	err = t.Execute(w, a.conf.Secrets["ROACH_DB"])
	if err != nil {
		log.Printf("ERROR: unable to execute the template: %v\n\n", err)
	}
}

func (a *App) dbInUse(hub *Hub, w http.ResponseWriter, r *http.Request) {
	c, err := a.u.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	client := &Client{hub: hub, conn: c, send: make(chan []byte, 256)}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}

func (a *App) getAllPlants(hub *Hub, w http.ResponseWriter, r *http.Request) {
	c, err := a.u.Upgrade(w, r, nil)
	if err != nil {
		log.Println("getAllPlants - upgrade: ", err)
		return
	}
	client := &Client{hub: hub, conn: c, send: make(chan []byte)}
	client.hub.register <- client

	go client.readPump()
	go client.writePump()
}
