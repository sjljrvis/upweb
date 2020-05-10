package helpers

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	sendgrid "github.com/sjljrvis/deploynow/lib/sendgrid"
)

func VerificationEmail(to string) {
	wd, _ := os.Getwd()
	t, err := template.ParseFiles(wd + "/templates/email_verify.html")
	if err != nil {
		fmt.Println(err)
	}

	data := struct {
		Email string
	}{
		Email: to,
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		fmt.Println(err)
	}

	result := tpl.String()
	fmt.Println(result)
	fmt.Println(data)
	sendgrid.SendEmail(to, "sejal@upweb.io", "Verify your upweb account", result)
}
