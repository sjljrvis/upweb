package lib

import (
	"fmt"
	"os"

	fs "github.com/sjljrvis/deploynow/lib/fs"
	"github.com/sjljrvis/deploynow/log"

	git "gopkg.in/src-d/go-git.v4"
)

// InitBare will initialize bare respository
func InitBare(path string) error {
	_, err := git.PlainInit(path, true)
	if err != nil {
		log.Info().Msgf("Error in Creating Bare repository")
		return err
	}
	fs.ChownR(path)
	return nil
}

// CreateHooks will initialize bare respository
func CreateHooks(path string) error {
	err := os.MkdirAll(path+"/hooks", os.FileMode(0775))
	if err != nil {
		return fmt.Errorf("error creating tabelspace folders: %v ", err.Error())
	}
	fs.ChownR(path + "/hooks")
	// err = fs.Copy("/Users/sejal/Projects/Personal/go/src/github.com/sjljrvis/deploynow/.githooks/pre-receive", path+"/hooks")
	// err = fs.Copy("/Users/sejal/Projects/Personal/go/src/github.com/sjljrvis/deploynow/.githooks/post-receive", path+"/hooks")
	return nil
}
