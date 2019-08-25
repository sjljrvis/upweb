package lib

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"
)

func CreateDir(path string) error {
	err := os.MkdirAll(path, os.FileMode(0775))
	if err != nil {
		return fmt.Errorf("error creating tabelspace folders: %v ", err.Error())
	}
	group, err := user.Lookup("www-data")
	if err != nil {
		return fmt.Errorf("error looking up postgres user user info")
	}
	uid, _ := strconv.Atoi(group.Uid)
	gid, _ := strconv.Atoi(group.Gid)

	err = syscall.Chown(path, uid, gid)
	return nil
}
