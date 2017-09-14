package wingedsnake

import (
	"errors"
	"net"
	"os"
	"sync"
)

const (
	globalKey = "WINGEDSNAKE_CONFIG"
)

var (
	errWindows = errors.New("not support windows")
)

// wingedSnake 实例
type wingedSnake struct {
	// 不可阻塞的初始化回调函数
	initCallback func(string, []net.Listener)
	// 阻塞的退出回调函数
	quitCallback func()
	running      []*os.Process
	wg           sync.WaitGroup
	iPID         int
}

// Main 运行实例
func Main(initFunc func(config string, socket []net.Listener), quitFunc func()) error {
	ws := &wingedSnake{
		initCallback: initFunc,
		quitCallback: quitFunc,
	}
	if len(os.Args) < 2 {
		return worker(ws)
	}
	return master(ws)
}
