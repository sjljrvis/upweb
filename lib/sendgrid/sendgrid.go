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

type Email struct {
	Value string `json:"email"`
}

type Personalization struct {
	To []Email `json:"to"`
}

type Content struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type SendGridPayload struct {
	Personalizations []Personalization `json:"personalizations"`
	From             Email             `json:"from"`
	Subject          string            `json:"subject"`
	Content          []Content         `json:"content"`
}

func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "\t")
	if err != nil {
		return in
	}
	return out.String()
}

func SendEmail(to, from, subject, body string) {

	_data := SendGridPayload{
		Personalizations: []Personalization{
			{
				To: []Email{
					{Value: to},
				},
			},
		},
		From:    Email{Value: from},
		Subject: subject,
		Content: []Content{
			{
				Type:  "text/html",
				Value: body,
			},
		},
	}

	payload, _ := json.Marshal(_data)
	fmt.Println(jsonPrettyPrint(string(payload)))

	url := "https://api.sendgrid.com/v3/mail/send"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("SENDGRID_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	defer resp.Body.Close()
	response_body, _ := ioutil.ReadAll(resp.Body)
	var response map[string]interface{}
	err = json.Unmarshal([]byte(response_body), &response)
	fmt.Println(response)
}
