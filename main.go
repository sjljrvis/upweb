package main

import (
	"net/http"

	"github.com/gorilla/handlers"
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
	corsObj := handlers.AllowedOrigins([]string{"*"})
	r.Use(router.LoggingMiddleware)
	http.Handle("/", handlers.CORS(corsObj)(r))
	http.ListenAndServe(":3000", nil)
}
