package helpers

import (
	"time"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"

)
var mySigningKey = []byte("captainjacksparrowsayshi")

// GenerateJWT will issue JWT token
func GenerateJWT(_id bson.ObjectId , email string ) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["_id"] = _id
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
			fmt.Errorf("Something Went Wrong: %s", err.Error())
			return "", err
	}

	return tokenString, nil
}
