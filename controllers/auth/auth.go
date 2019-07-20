package controller

import (
	"encoding/json"
	"net/http"

	// "github.com/gorilla/mux"
	Helper "github.com/sjljrvis/deploynow/helpers"
	UserModel "github.com/sjljrvis/deploynow/models/user"
	"gopkg.in/mgo.v2/bson"
)

type auth struct {
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"`
}

// Login controller
func Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var query map[string]interface{}
	query = make(map[string]interface{})
	var authData auth

	if err := json.NewDecoder(r.Body).Decode(&authData); err != nil {
		Helper.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	query["email"] = authData.Email
	result, err := UserModel.FindOne(query)

	if err != nil {
		Helper.RespondWithError(w, 200, "email or password is")
		return
	}

	hashCheck := Helper.CheckPasswordHash(authData.Password, result.Password)
	if hashCheck {
		token, _ := Helper.GenerateJWT(result.ID, result.Email)
		Helper.RespondWithJSON(w, 200, map[string]string{"message": "login success", "token": token, "email": result.Email})
	} else {
		Helper.RespondWithError(w, 200, "Some error")
	}
}

// Register controller
func Register(w http.ResponseWriter, r *http.Request) {
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
