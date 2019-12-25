package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sjljrvis/deploynow/log"
)

type accessToken struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	// TokenType   string `json:"token_type"`
}

func AccessToken(code, state string) {
	var accessTokenData accessToken
	data := map[string]interface{}{
		"code":          code,
		"state":         state,
		"client_id":     "821083a95a975c302f45",
		"client_secret": "11bb8e2adeff0efc5c884a52071355ccd894760e",
		"redirect_uri":  "http://localhost:3001/#/oauth",
	}
	url := "https://github.com/login/oauth/access_token"
	payload, err := json.Marshal(data)
	log.Info().Msg("Calling github api--1")
	if err != nil {
		fmt.Println(err)
	}
	response, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	defer response.Body.Close()
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	err = json.Unmarshal([]byte(body), &accessTokenData)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(accessTokenData)
	fmt.Println(string(body))
}
