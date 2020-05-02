package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

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
		"client_id":     os.Getenv("GITHUB_CLIENT_ID"),
		"client_secret": os.Getenv("GITHUB_CLIENT_SECRET"),
		"redirect_uri":  os.Getenv("GITHUB_REDIRECT_URL"),
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

func Repositories(user string) []map[string]interface{} {
	api_url := fmt.Sprintf("https://api.github.com/users/%s/repos?per_page=50", user)
	req, err := http.NewRequest("GET", api_url, nil)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var dt []map[string]interface{}
	_ = json.Unmarshal(body, &dt)
	return dt
}
