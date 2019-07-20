package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	Helper "github.com/sjljrvis/deploynow/helpers"
	UserModel "github.com/sjljrvis/deploynow/models/user"
	"gopkg.in/mgo.v2/bson"
)

// GetAll controller
func GetAll(w http.ResponseWriter, r *http.Request) {
	var results []UserModel.User
	results, err := UserModel.FindAll()
	if err != nil {
		Helper.RespondWithError(w, 200, err.Error())
	} else {
		Helper.RespondWithJSON(w, 200, results)
	}
}

// GetOne controller
func GetOne(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	result, err := UserModel.FindOneByID(params["id"])
	if err != nil {
		Helper.RespondWithError(w, 200, err.Error())
	} else {
		Helper.RespondWithJSON(w, 200, result)
	}
}

// Create controller
func Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var user UserModel.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		Helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	user.ID = bson.NewObjectId()
	hashed, _ := Helper.HashPassword(user.Password)
	user.Password = hashed
	err := UserModel.Create(user)
	if err != nil {
		Helper.RespondWithError(w, 200, err.Error())
	} else {
		Helper.RespondWithJSON(w, 200, map[string]string{"message": "User Created successfully"})
	}
}

// Search controller
func Search(w http.ResponseWriter, r *http.Request) {

	var query map[string]interface{}
	query = make(map[string]interface{})

	keys := r.URL.Query()

	for item := range keys {
		query[item] = keys[item][0]
	}

	result, err := UserModel.FindByQuery(query)
	if err != nil {
		Helper.RespondWithError(w, 200, err.Error())
	} else {
		Helper.RespondWithJSON(w, 200, result)
	}
}
