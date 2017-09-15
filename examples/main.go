package main

import (
	log "github.com/928799934/log4go.v1"
	ws "github.com/928799934/wingedsnake"
	"net"
)

func main() {
	if err := ws.Main(InitFunc, QuitFunc); err != nil {
		panic(err)
	}
}

// InitFunc 测试
func InitFunc(config string, listeners []net.Listener) {
	log.LoadConfiguration(config)
	startHTTPListen(listeners)
	log.Info("init")
}

// QuitFunc 测试
func QuitFunc() {
	stopHTTPListen()
	log.Info("quit")
	log.Close()
}
