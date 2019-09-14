package lib

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	Helper "github.com/sjljrvis/deploynow/helpers"
)

func getConfig(name string, port string) string {
	conf := `
	server {
		listen 80; 
		server_name ` + name + `.tocstack.com;
		location / {
		 proxy_set_header X-Real-IP $remote_addr;
		 proxy_set_header Host $host;
		 proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
		 proxy_pass http://localhost:` + port + `;
	 }
	}`
	return conf
}

func symlink(name string) error {
	err := os.Symlink(name, "sites_enabled/name")
	if err != nil {
		return err
	}
	return nil
}

//Reload reloads nginx
func Reload() {
	cmd := exec.Command("sudo", "service", "nginx", "reload")
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}

// WriteConfig Creates ginx config
func WriteConfig(name string, port string) {
	conf := []byte(getConfig(name, port))
	err := ioutil.WriteFile(name+".nginx.conf", conf, 0777)
	if err != nil {
		fmt.Println(err)
	}
}

//Writehtpasswd secures directory with http password
func Writehtpasswd(path, username, password string) error {
	text := username + `:` + Helper.GetMD5Hash(password)
	err := ioutil.WriteFile(path, []byte(text), 0777)
	if err != nil {
		return err
	}
	return nil
}
