package main

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/sjljrvis/deploynow/log"
	"github.com/sjljrvis/deploynow/router"
	DB "github.com/sjljrvis/deploynow/db"
)

func init() {
	log.Info().Msgf("Starting server at port 3000")
}

func main() {
	DB.Init()
	r := router.NewRouter()

	corsObj := handlers.AllowedOrigins([]string{"*"})

	http.Handle("/", handlers.CORS(corsObj)(r))
	http.ListenAndServe(":3000", nil)
}
