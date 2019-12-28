package lib

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/sjljrvis/deploynow/log"
)

type accessToken struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func AccessToken(code, state string) accessToken {
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
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal([]byte(body), &accessTokenData)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	return accessTokenData
}

func Profile(token string) []byte {
	url := "https://api.github.com/user?access_token=" + token

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body
}
