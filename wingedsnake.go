package wingedsnake

import (
	"errors"
	"net"
	"os"
	"sync"
)

const (
	globalKey   = "WINGEDSNAKE_CONFIG"
	affinityKey = "WINGEDSNAKE_AFFINITY"
)

var (
	errWindows = errors.New("not support windows")
)

// wingedSnake 实例
type wingedSnake struct {
	// 不可阻塞的初始化回调函数
	initCallback func(string, []net.Listener, []net.Listener, []net.Listener)
	// 阻塞的退出回调函数
	quitCallback func()
	closeEvent   chan bool
	list         []*os.Process
	wg           sync.WaitGroup
	pidFilePath  string
	iPID         int
}

// Main 运行实例
func Main(initFunc func(config string, socket, ping, pprof []net.Listener), quitFunc func()) error {
	// 检测系统
	if err := supportOS(); err != nil {
		return err
	}
	ws := &wingedSnake{
		initCallback: initFunc,
		quitCallback: quitFunc,
		closeEvent:   make(chan bool, 0),
	}
	if len(os.Args) < 2 {
		return worker(ws)
	}
	return master(ws)
}
