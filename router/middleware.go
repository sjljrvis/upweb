package router

import (
	"context"
	"net/http"
	"strings"
	"time"

	. "github.com/sjljrvis/deploynow/db"
	Helper "github.com/sjljrvis/deploynow/helpers"
	"github.com/sjljrvis/deploynow/log"

	"github.com/sjljrvis/deploynow/models"

	"github.com/dgrijalva/jwt-go"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		latency := time.Since(start)
		ms := float32(latency / time.Millisecond)
		log.Info().Msgf("time=%.3fms url=%s IP=%s", ms, r.RequestURI, r.RemoteAddr)
	})
}

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
					claims := token.Claims.(jwt.MapClaims)
					user := models.User{}
					err := DB.Where("uuid = ?", claims["uuid"]).First(&user).Error
					if err != nil {
						Helper.RespondWithError(w, 403, err.Error())
						return
					}
					ctx := context.WithValue(req.Context(), "user", user)
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
