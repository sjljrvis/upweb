package main

import (
	"net/http"

	DB "github.com/sjljrvis/deploynow/db"

	"github.com/sjljrvis/deploynow/log"
	"github.com/sjljrvis/deploynow/router"
	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load()
	log.Info().Msgf("Starting server at port 3000")
}

func main() {

	DB.Init()
	r := router.NewRouter()
	r.Use(router.LoggingMiddleware)

	http.Handle("/", router.CorsMiddleware(r))
	http.ListenAndServe(":3000", nil)
}
