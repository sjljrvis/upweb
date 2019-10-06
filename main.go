package main

import (
	"net/http"

	"github.com/gorilla/handlers"
	DB "github.com/sjljrvis/deploynow/db"
	do "github.com/sjljrvis/deploynow/lib/digitalocean"
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
	do.CreateDNS("sejal")
	r := router.NewRouter()
	corsObj := handlers.AllowedOrigins([]string{"*"})

	http.Handle("/", handlers.CORS(corsObj)(r))
	http.ListenAndServe(":3000", nil)
}
