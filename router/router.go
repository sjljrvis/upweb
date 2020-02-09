package router

import (
	"net/http"

	"github.com/gorilla/mux"

	ActivityController "github.com/sjljrvis/deploynow/controllers/activity"
	AuthController "github.com/sjljrvis/deploynow/controllers/auth"
	BuildsController "github.com/sjljrvis/deploynow/controllers/builds"
	ContainerController "github.com/sjljrvis/deploynow/controllers/containers"
	RepositoryController "github.com/sjljrvis/deploynow/controllers/repository"
	VariableController "github.com/sjljrvis/deploynow/controllers/variable"

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
	authRouter.HandleFunc("/github/oauth", AuthMiddleware(AuthController.GithubAuth)).Methods("POST")

	/*
		user subrouter
		handle  REST-api /user here
	*/

	userRouter := r.PathPrefix("/api/v1/user").Subrouter()
	userRouter.HandleFunc("/password", AuthMiddleware(UserController.UpdatePassword)).Methods("PUT")
	userRouter.HandleFunc("/", UserController.GetAll).Methods("GET")
	userRouter.HandleFunc("/", UserController.Create).Methods("POST")
	userRouter.HandleFunc("/{id}", UserController.Delete).Methods("DELETE")
	userRouter.HandleFunc("/{id}", AuthMiddleware(UserController.Update)).Methods("PUT")
	userRouter.HandleFunc("/{uuid}/github", UserController.GetGithubAccount).Methods("GET")
	userRouter.HandleFunc("/{uuid}/github", UserController.RemoveGithubAccount).Methods("DELETE")

	// userRouter.HandleFunc("/search/", UserController.Search).Methods("GET")

	/*
		repository subrouter
		handle  REST-api /user here
	*/

	repositoryRouter := r.PathPrefix("/api/v1/repository").Subrouter()
	// repositoryRouter.Use(AuthMiddleware)
	repositoryRouter.HandleFunc("/", AuthMiddleware(RepositoryController.GetAll)).Methods("GET")
	repositoryRouter.HandleFunc("/github/", AuthMiddleware(RepositoryController.GetGithubRepos)).Methods("GET")
	repositoryRouter.HandleFunc("/{uuid}", AuthMiddleware(RepositoryController.Get)).Methods("GET")
	repositoryRouter.HandleFunc("/", AuthMiddleware(RepositoryController.Create)).Methods("POST")
	repositoryRouter.HandleFunc("/{id}", AuthMiddleware(RepositoryController.Update)).Methods("PUT")
	repositoryRouter.HandleFunc("/{id}", AuthMiddleware(RepositoryController.Delete)).Methods("DELETE")
	/*
		repository subrouter
		handle  REST-api /user here
	*/

	variableRouter := r.PathPrefix("/api/v1/variable").Subrouter()
	// repositoryRouter.Use(AuthMiddleware)
	variableRouter.HandleFunc("/{repository_id}", AuthMiddleware(VariableController.GetAll)).Methods("GET")
	variableRouter.HandleFunc("/{repository_id}/{id}", AuthMiddleware(VariableController.Get)).Methods("GET")
	variableRouter.HandleFunc("/{repository_id}", AuthMiddleware(VariableController.Create)).Methods("POST")
	// variableRouter.HandleFunc("/{repository_id}/{id}", AuthMiddleware(RepositoryController.Update)).Methods("PUT")
	// variableRouter.HandleFunc("/{repository_id}/{id}", AuthMiddleware(RepositoryController.Delete)).Methods("DELETE")

	activityRouter := r.PathPrefix("/api/v1/activity").Subrouter()
	// repositoryRouter.Use(AuthMiddleware)
	activityRouter.HandleFunc("/{repository_id}", AuthMiddleware(ActivityController.GetAll)).Methods("GET")

	/*
		repository subrouter
		handle  REST-api /user here
	*/

	buildsRouter := r.PathPrefix("/api/v1/build").Subrouter()
	buildsRouter.HandleFunc("/{repository_id}", AuthMiddleware(BuildsController.GetAll)).Methods("GET")
	// variableRouter.HandleFunc("/{repository_id}/{id}", AuthMiddleware(RepositoryController.Update)).Methods("PUT")
	// variableRouter.HandleFunc("/{repository_id}/{id}", AuthMiddleware(RepositoryController.Delete)).Methods("DELETE")

	/*
		container subrouter
		handle  REST-api /user here
	*/

	containerRouter := r.PathPrefix("/api/v1/container").Subrouter()
	// repositoryRouter.Use(AuthMiddleware)
	containerRouter.HandleFunc("/build", ContainerController.Build).Methods("POST")
	containerRouter.HandleFunc("/logs/{uuid}", ContainerController.Getlogs).Methods("GET")
	containerRouter.HandleFunc("/dev", ContainerController.BuildLogs).Methods("GET")

	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", fs))
	return r
}
