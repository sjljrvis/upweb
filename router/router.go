package router

import (
	"net/http"

	"github.com/gorilla/mux"

	UserController "github.com/sjljrvis/deploynow/controllers/user"
	AuthController "github.com/sjljrvis/deploynow/controllers/auth"

)

// NewRouter is router pointer
func NewRouter() *mux.Router {
	fs := http.FileServer(http.Dir("./public/"))
	r := mux.NewRouter()


	/*
		user subrouter
		handle  REST-api /user here
	*/

	authRouter := r.PathPrefix("/api/v1/auth").Subrouter()
	authRouter.HandleFunc("/login", AuthController.Login).Methods("POST")
	authRouter.HandleFunc("/register", AuthController.Register).Methods("POST")


	/*
		user subrouter
		handle  REST-api /user here
	*/

	userRouter := r.PathPrefix("/api/v1/user").Subrouter()
	userRouter.HandleFunc("/", UserController.GetAll).Methods("GET")
	userRouter.HandleFunc("/{id}", UserController.GetOne).Methods("GET")
	userRouter.HandleFunc("/", UserController.Create).Methods("POST")
	userRouter.HandleFunc("/search/", UserController.Search).Methods("GET")

	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", fs))
	return r
}
