package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/foomo/htpasswd"
	DB "github.com/sjljrvis/deploynow/db"

	"github.com/sjljrvis/deploynow/log"
	"github.com/sjljrvis/deploynow/router"
	"github.com/subosito/gotenv"
)

const PrefixCryptApr1 = "$apr1$"

func init() {
	env_path := fmt.Sprintf("./configs/.%s.env", os.Getenv("ENV"))
	gotenv.Load(env_path)
	log.Info().Msgf("environment %s", os.Getenv("ENV"))
	log.Info().Msgf("Starting server at port %s", os.Getenv("PORT"))

	passwords := htpasswd.HashedPasswords(map[string]string{})
	passwords.SetPassword("sjl", "123456", htpasswd.HashAPR1)
	fmt.Println(passwords)
}

func main() {

	DB.Init()
	r := router.NewRouter()
	r.Use(router.LoggingMiddleware)

	http.Handle("/", router.CorsMiddleware(r))
	http.ListenAndServe(":3000", nil)
}
