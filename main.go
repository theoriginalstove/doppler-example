package main

import (
	"flag"
	"os"

	"gitlab.com/steven.t/doppler-example/app"
)

func main() {
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":5000"
	}
	flag.Parse()
	a := app.Configure("better_secrets", app.WithAddr(addr))
	a.Server.ListenAndServe()
}
