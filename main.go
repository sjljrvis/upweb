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

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers:", "Origin, Content-Type, X-Auth-Token, Authorization")
			w.Header().Set("Content-Type", "application/json")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {

	DB.Init()
	r := router.NewRouter()
	r.Use(router.LoggingMiddleware)

	http.Handle("/", corsMiddleware(r))
	http.ListenAndServe(":3000", nil)
}
