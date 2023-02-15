package app

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"gitlab.com/steven.t/doppler-example/config"
	"gitlab.com/steven.t/doppler-example/db"
)

type App struct {
	name   string
	Db     *db.Instance
	conf   *config.Config
	Server Server
	u      *websocket.Upgrader
	hub    *Hub
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
	err := a.conf.GetDopplerSecrets()
	if err != nil {
		http.Error(w, fmt.Sprintf("%d - Server Error - unable to refresh secrets", http.StatusInternalServerError), http.StatusInternalServerError)
	}

	db := a.conf.Secrets["ROACH_DB"]
	a.hub.broadcast <- []byte(db)

	_, err = w.Write([]byte("ok"))
	if err != nil {
		log.Println("err writing to http response")
	}
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

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, []byte{}, -1))
		c.hub.broadcast <- message

	}
}

const (
	writeWait        = 10 * time.Second
	maxMessageSize   = 8192
	pongWait         = 60 * time.Second
	pingPeriod       = (pongWait * 9) / 10
	closeGracePeriod = 10 * time.Second
)

func (a *App) messageIn(ws *websocket.Conn, w io.Writer) {
	defer ws.Close()
	ws.SetReadLimit(maxMessageSize)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		if _, err := w.Write(message); err != nil {
			break
		}
	}
}

func ping(ws *websocket.Conn, done chan struct{}) {
	ticker := time.NewTimer(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				log.Println("ping:", err)
			}
		case <-done:
			return
		}
	}
}

type Hub struct {
	clients map[*Client]struct{}

	broadcast chan []byte

	register chan *Client

	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]struct{}),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = struct{}{}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}
