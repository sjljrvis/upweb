package router

import (
	"net/http"

	"github.com/gorilla/mux"

	AuthController "github.com/sjljrvis/deploynow/controllers/auth"
	RepositoryController "github.com/sjljrvis/deploynow/controllers/repository"
	UserController "github.com/sjljrvis/deploynow/controllers/user"
	Helper "github.com/sjljrvis/deploynow/helpers"
)

func testEndpoint(w http.ResponseWriter, req *http.Request) {
	Helper.RespondWithJSON(w, 200, map[string]string{"message": "User Created successfully"})
}

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

	/*
		repository subrouter
		handle  REST-api /user here
	*/

	repositoryRouter := r.PathPrefix("/api/v1/repository").Subrouter()
	// repositoryRouter.Use(AuthMiddleware)
	repositoryRouter.HandleFunc("/", AuthMiddleware(RepositoryController.GetAll)).Methods("GET")
	repositoryRouter.HandleFunc("/{id}", RepositoryController.GetOne).Methods("GET")
	repositoryRouter.HandleFunc("/", RepositoryController.Create).Methods("POST")
	repositoryRouter.HandleFunc("/search/", RepositoryController.Search).Methods("GET")

	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", fs))
	return r
}
