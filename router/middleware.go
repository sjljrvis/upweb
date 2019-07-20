package router

import (
	"context"
	"net/http"
	"strings"

	Helper "github.com/sjljrvis/deploynow/helpers"

	"github.com/dgrijalva/jwt-go"
)

// AuthMiddleware will check headers
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					return []byte("captainjacksparrowsayshi"), nil
				})
				if error != nil {
					Helper.RespondWithError(w, 403, error.Error())
					return
				}
				if token.Valid {
					ctx := context.WithValue(req.Context(), "decoded", token.Claims)
					next(w, req.WithContext(ctx))
				} else {
					Helper.RespondWithError(w, 403, "Authorization Invalid")
				}
			}
		} else {
			Helper.RespondWithError(w, 403, "An authorization header is required")
		}
	})
}
