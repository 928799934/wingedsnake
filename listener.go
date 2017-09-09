package wingedsnake

import (
	"net"
	"os"
	"strings"
)

// getListenersByFD 根据 输入fd 进行监听
func getListenersByFD(pprof, socket, ping []string) ([]net.Listener, []net.Listener, []net.Listener, error) {
	var (
		pprofListeners  []net.Listener
		socketListeners []net.Listener
		pingListeners   []net.Listener
	)
	// 0:标准输入
	// 1:标准输出
	// 2:标准错误输出
	fd := uintptr(3)
	for num := len(pprof); num > 0; num-- {
		f := os.NewFile(fd, "pprof listen")
		l, err := net.FileListener(f)
		if err != nil {
			logf("net.FileListener(file) error(%v)", err.Error())
			return nil, nil, nil, err
		}
		pprofListeners = append(pprofListeners, l)
		fd++
	}
	for num := len(socket); num > 0; num-- {
		f := os.NewFile(fd, "socket listen")
		l, err := net.FileListener(f)
		if err != nil {
			logf("net.FileListener(file) error(%v)", err.Error())
			return nil, nil, nil, err
		}
		socketListeners = append(socketListeners, l)
		fd++
	}
	for num := len(ping); num > 0; num-- {
		f := os.NewFile(fd, "ping listen")
		l, err := net.FileListener(f)
		if err != nil {
			logf("net.FileListener(file) error(%v)", err.Error())
			return nil, nil, nil, err
		}
		pingListeners = append(pingListeners, l)
		fd++
	}
	return pprofListeners, socketListeners, pingListeners, nil
}

// getListeners 根据地址 监听
func getListeners(pprof, socket, ping []string) ([]net.Listener, []net.Listener, []net.Listener, error) {
	var (
		l               net.Listener
		err             error
		pprofListeners  []net.Listener
		socketListeners []net.Listener
		pingListeners   []net.Listener
	)

	for _, v := range pprof {
		if strings.Index(v, ":") == -1 {
			l, err = net.Listen("unix", v)
			os.Chmod(v, 0777)
		} else {
			l, err = net.Listen("tcp", v)
		}
		if err != nil {
			return nil, nil, nil, err
		}
		pprofListeners = append(pprofListeners, l)
	}

	for _, v := range socket {
		if strings.Index(v, ":") == -1 {
			l, err = net.Listen("unix", v)
			os.Chmod(v, 0777)
		} else {
			l, err = net.Listen("tcp", v)
		}
		if err != nil {
			return nil, nil, nil, err
		}
		socketListeners = append(socketListeners, l)
	}

	for _, v := range ping {
		if strings.Index(v, ":") == -1 {
			l, err = net.Listen("unix", v)
			os.Chmod(v, 0777)
		} else {
			l, err = net.Listen("tcp", v)
		}
		if err != nil {
			return nil, nil, nil, err
		}
		pingListeners = append(pingListeners, l)
	}

	return pprofListeners, socketListeners, pingListeners, nil
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
