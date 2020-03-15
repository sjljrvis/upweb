package lib

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/sjljrvis/deploynow/log"
)

func getConfig(name string, port string) string {
	conf := `
	server {
    server_name ` + name + `.upweb.io;

    listen 443 ssl; # managed by Certbot
    ssl_certificate /etc/letsencrypt/live/upweb.io/fullchain.pem; # managed by Certbot
    ssl_certificate_key /etc/letsencrypt/live/upweb.io/privkey.pem; # managed by Certbot
    include /etc/letsencrypt/options-ssl-nginx.conf; # managed by Certbot
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem; # managed by Certbot

    location / {
             proxy_set_header X-Real-IP $remote_addr;
             proxy_set_header Host $host;
             proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
						 proxy_pass http://localhost:` + port + `;
          }
		}

	server {
		listen 80;
    listen [::]:80;
    server_name ` + name + `.upweb.io;
    return 301 https://$host$request_uri;
    return 404; # managed by Certbot
	}`
	return conf
}

func Symlink(name string) error {
	confPath := path.Join(os.Getenv("NGINX_SITES_AVAILABLE"), name+".upweb.io")
	sitesEnabled := path.Join(os.Getenv("NGINX_SITES_ENABLED"), name+".upweb.io")
	err := os.Symlink(confPath, sitesEnabled)
	if err != nil {
		return err
	}
	return nil
}

//Reload reloads nginx
func Reload() error {
	cmd := exec.Command("sudo", "service", "nginx", "reload")
	err := cmd.Run()
	if err != nil {
		log.Error().Err(err)
		return err
	}
	return nil
}

// WriteConfig Creates ginx config
func WriteConfig(name string, port string) {
	confPath := path.Join(os.Getenv("NGINX_SITES_AVAILABLE"), name+".upweb.io")
	conf := []byte(getConfig(name, port))
	err := ioutil.WriteFile(confPath, conf, 0777)
	if err != nil {
		fmt.Println(err)
	}
}

//Writehtpasswd secures directory with http password
func Writehtpasswd(username, password string) error {
	_path := path.Join(os.Getenv("ROOT_DIR"), username, "htpasswd")
	text := username + `:` + password
	err := ioutil.WriteFile(_path, []byte(text), 0777)
	if err != nil {
		return err
	}
	return nil
}
