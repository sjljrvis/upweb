package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var domain = "tocstack.com"

func CreateDNS(repository_name string) error {
	data := map[string]interface{}{
		"type":     "A",
		"name":     repository_name,
		"data":     "139.59.69.254",
		"priority": nil,
		"port":     nil,
		"ttl":      1800,
		"weight":   nil,
		"flags":    nil,
		"tag":      nil,
	}
	url := os.Getenv("DIGITAL_OCEAN_HOST") + "/" + domain + "/records"
	payload, _ := json.Marshal(data)
	response, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	defer response.Body.Close()
	if err != nil {
		return err
	}
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Sprint(body)
	return nil
}

func GetDNS() {

}

func RemoveDNS() {

}
