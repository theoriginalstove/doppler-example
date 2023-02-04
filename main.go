package main

import (
	"log"
	"net/http"

	"gitlab.com/steven.t/doppler-example/app"
)

func main() {
	app := app.Configure("")
	log.Printf("listening on port %s", "5005")
	http.ListenAndServe(":5005", app.Router)
}
