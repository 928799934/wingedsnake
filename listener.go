// +build dragonfly freebsd linux netbsd openbsd solaris

package wingedsnake

import (
	"net"
	"os"
	"strings"
)

// getListenersByFD 根据 输入fd 进行监听
func getListenersByFD(socket []string) ([]net.Listener, error) {
	var (
		list []net.Listener
	)
	// 0:标准输入
	// 1:标准输出
	// 2:标准错误输出
	fd := uintptr(3)
	for num := len(socket); num > 0; num-- {
		f := os.NewFile(fd, "socket listen")
		l, err := net.FileListener(f)
		if err != nil {
			logf("net.FileListener(file) error(%v)", err.Error())
			return nil, err
		}
		list = append(list, l)
		fd++
	}
	return list, nil
}

// getListeners 根据地址 监听
func getListeners(addrs []string) ([]net.Listener, error) {
	var list []net.Listener

	for _, v := range addrs {
		if strings.Index(v, ":") != -1 {
			l, err := net.Listen("tcp", v)
			if err != nil {
				return nil, err
			}
			list = append(list, l)
			continue
		}
		l, err := net.Listen("unix", v)
		if err != nil {
			return nil, err
		}
		os.Chmod(v, 0777)
		list = append(list, l)
	}
	return list, nil
}

// getFileByListener 获取file
func getFileByListener(l net.Listener) *os.File {
	switch inst := l.(type) {
	case *net.TCPListener:
		f, err := inst.File()
		if err != nil {
			logf("inst.File() error(%v)", err)
			return nil
		}
		return f
	case *net.UnixListener:
		f, err := inst.File()
		if err != nil {
			logf("inst.File() error(%v)", err)
			return nil
		}
		return f
	}
	return nil
}

func closeListeners(list []net.Listener) {
	for _, l := range list {
		switch inst := l.(type) {
		case *net.TCPListener:
			inst.Close()
		case *net.UnixListener:
			inst.SetUnlinkOnClose(true)
			inst.Close()
		}
	}
}
