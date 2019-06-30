package main

import (
	"net/http"
	"github.com/sjljrvis/deploynow/log"
	"github.com/gorilla/handlers"
	"github.com/sjljrvis/deploynow/router"
	tool "github.com/sjljrvis/deploynow/tools"
)

func init() {
	log.Info().Msgf("Starting server at port 3000")
}

func main() {
	tool.MongoConnect("mongodb://sejal:sejal@ds117935.mlab.com:17935/mobystore", "mobystore")
	r := router.NewRouter()

	corsObj := handlers.AllowedOrigins([]string{"*"})

	http.Handle("/", handlers.CORS(corsObj)(r))
	http.ListenAndServe(":3000", nil)
}
