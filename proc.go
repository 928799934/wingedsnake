// +build dragonfly freebsd linux netbsd openbsd solaris

package wingedsnake

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

func createPidFile(ws *wingedSnake, pidPath string) error {
	szPath := []byte(pidPath)
	if szPath[len(szPath)-1] != '/' {
		pidPath = pidPath + "/"
	}
	_, file := filepath.Split(os.Args[0])
	pidFilePath := pidPath + file + ".pid"
	ws.iPID = os.Getpid()
	strPID := strconv.Itoa(ws.iPID)
	if err := ioutil.WriteFile(pidFilePath, []byte(strPID), 0644); err != nil {
		logf("ioutil.WriteFile(%v,%v,0644) error(%v)", pidFilePath, strPID, err)
		return err
	}
	return nil
}

func removePidFile(ws *wingedSnake, pidPath string) error {
	szPath := []byte(pidPath)
	if szPath[len(szPath)-1] != '/' {
		pidPath = pidPath + "/"
	}
	_, file := filepath.Split(os.Args[0])
	pidFilePath := pidPath + file + ".pid"
	buf, err := ioutil.ReadFile(pidFilePath)
	if err != nil {
		logf("ioutil.ReadFile(%v) error(%v)", pidFilePath, err)
		return err
	}
	strPID := string(buf)
	pid, err := strconv.Atoi(strPID)
	if err != nil {
		logf("strconv.Atoi(%v) error(%v)", strPID, err)
		return err
	}
	if pid != ws.iPID {
		return nil
	}
	if err := os.Remove(pidFilePath); err != nil {
		logf("os.Remove(%v) error(%v)", pidFilePath, err)
		return err
	}
	return nil
}
