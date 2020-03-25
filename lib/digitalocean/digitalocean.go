package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/sjljrvis/deploynow/log"
)

var domain = "upweb.io"

var i interface{}

type digitalOcean struct {
	Domain string
	IP     string
	Token  string
}

// CreateDNS will create a A name record in cloud provide
func CreateDNS(repository_name string) (string, error) {

	doClient := new(digitalOcean)
	doClient.Domain = domain
	doClient.IP = "167.71.237.137"
	doClient.Token = "38936697bd54da1c86dbf68e737f49cd60492d5a8c31d7ce4b6b76bce1450b06"

	data := map[string]interface{}{
		"type":     "A",
		"name":     repository_name,
		"data":     doClient.IP,
		"priority": nil,
		"port":     nil,
		"ttl":      1800,
		"weight":   nil,
		"flags":    nil,
		"tag":      nil,
	}
	payload, err := json.Marshal(data)

	url := fmt.Sprintf("%s/%s/records", os.Getenv("DIGITAL_OCEAN_HOST"), domain)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+doClient.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var response map[string]interface{}
	err = json.Unmarshal([]byte(body), &response)
	response_map := response["domain_record"].(map[string]interface{})
	_id := response_map["id"].(float64)
	return strconv.FormatFloat(_id, 'f', 6, 64), err
}

// GetDNS will fetch all records
func GetDNS(record_id string) (interface{}, error) {
	doClient := new(digitalOcean)
	doClient.Domain = domain
	doClient.IP = "167.71.237.137"
	doClient.Token = "38936697bd54da1c86dbf68e737f49cd60492d5a8c31d7ce4b6b76bce1450b06"

	url := fmt.Sprintf("%s/%s/records/%s", os.Getenv("DIGITAL_OCEAN_HOST"), domain, record_id)
	req, err := http.NewRequest("GET", url, nil)

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+doClient.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal([]byte(body), &i)
	log.Info().Msg(string(body))
	return i, err
}

// RemoveDNS will remove A name record fron cloud Provider
func RemoveDNS(record_id string) (interface{}, error) {
	doClient := new(digitalOcean)
	doClient.Domain = domain
	doClient.IP = "167.71.237.137"
	doClient.Token = "38936697bd54da1c86dbf68e737f49cd60492d5a8c31d7ce4b6b76bce1450b06"

	url := fmt.Sprintf("%s/%s/records/%s", os.Getenv("DIGITAL_OCEAN_HOST"), domain, record_id)
	log.Info().Msg(url)
	req, err := http.NewRequest("DELETE", url, nil)

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+doClient.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal([]byte(body), &i)
	log.Info().Msg(string(body))
	return i, err
}
