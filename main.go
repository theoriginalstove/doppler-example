package main

import (
	"os"

	"gitlab.com/steven.t/doppler-example/app"
)

func main() {
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":5005"
	}
	a := app.Configure("", app.WithAddr(addr))
	a.Server.ListenAndServe()
}
