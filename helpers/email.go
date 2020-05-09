package helpers

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	sendgrid "github.com/sjljrvis/deploynow/lib/sendgrid"
)

func VerificationEmail() {
	wd, _ := os.Getwd()
	t, err := template.ParseFiles(wd + "/templates/email_verify.html")
	if err != nil {
		fmt.Println(err)
	}

	key := "some strings"

	data := struct {
		Key string
	}{
		Key: key,
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		fmt.Println(err)
	}

	result := tpl.String()
	fmt.Println(result)
	fmt.Println(data)
	sendgrid.SendEmail("sjlchougule@gmail.com", "sejal@upweb.io", "Verify your upweb account", result)

	// fmt.Println(res)

	// var accessTokenData accessToken
	// _data = map[string]interface{}{
	// 	"personalizations": ["sejal"]
	// }
	// url := "https://github.com/login/oauth/access_token"

	// payload, err := json.Marshal(data)
	// req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	// req.Header.Set("Accept", "application/json")
	// req.Header.Set("Content-Type", "application/json")

	// client := &http.Client{}
	// resp, err := client.Do(req)
	// if err != nil {
	// 	log.Error().Msg(err.Error())
	// }
	// defer resp.Body.Close()
	// body, _ := ioutil.ReadAll(resp.Body)
	// err = json.Unmarshal([]byte(body), &accessTokenData)
	// if err != nil {
	// 	log.Error().Msg(err.Error())
	// }
	// return accessTokenData
}

// curl --request POST \
// --url https://api.sendgrid.com/v3/mail/send \
// --header "Authorization: Bearer $SENDGRID_API_KEY" \
// --header 'Content-Type: application/json' \
// --data '
// {"personalizations": [
// 	{"to":
// 	[
// 		{"email": "sjljarvis@gmail.com"}
// 		]
// 		}],"from": {"email": "sejal@upweb.io"},"subject": "Sending with SendGrid is Fun","content": [{"type": "text/plain", "value": "and easy to do anywhere, even with cURL"}]}'
