package lib

import (
	"fmt"
	"os"
	"path"

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
func CreateHooks(dir string) error {
	err := os.MkdirAll(dir+"/hooks", os.FileMode(0775))
	if err != nil {
		return fmt.Errorf("error creating tabelspace folders: %v ", err.Error())
	}
	fs.ChownR(dir + "/hooks")
	err = fs.Copy(path.Join(os.Getenv("PROJECT_DIR"), ".githooks", "pre-receive"), path.Join(dir, "hooks", "pre-receive"))
	err = fs.Copy(path.Join(os.Getenv("PROJECT_DIR"), ".githooks", "post-receive"), path.Join(dir, "hooks", "post-receive"))
	return nil
}
