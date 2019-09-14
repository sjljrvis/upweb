package lib

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path"
	"runtime"
	"strconv"
	"syscall"

	"github.com/sjljrvis/deploynow/log"
)

//CreateDir will create dir with group
func CreateDir(path string) error {
	err := os.MkdirAll(path, os.FileMode(0775))
	if err != nil {
		return fmt.Errorf("error creating tabelspace folders: %v ", err.Error())
	}
	ChownR(path)
	return nil
}

func Chown(path, group string) {
	if runtime.GOOS != "darwin" {
		group, err := user.Lookup(group)
		if err != nil {
			log.Info().Msgf(err.Error())
		}
		uid, _ := strconv.Atoi(group.Uid)
		gid, _ := strconv.Atoi(group.Gid)

		err = syscall.Chown(path, uid, gid)
	}
}

func ChownR(path string) {
	if runtime.GOOS != "darwin" {
		cmd := exec.Command("chown www-data:www-data", "-R", path)
		cmd.Run()
	}
}

// Copy file from src to destination
func Copy(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

// DirCopy file from src to destination
func DirCopy(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = DirCopy(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = DirCopy(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}
