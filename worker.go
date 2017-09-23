// +build dragonfly freebsd linux netbsd openbsd solaris

package wingedsnake

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// worker 工作进程函数
func worker(ws *wingedSnake) error {
	// 读取环境变量中的配置
	jsonData := []byte(os.Getenv(globalKey))
	if len(jsonData) == 0 {
		return fmt.Errorf("lost config path")
	}

	conf := &config{}
	if err := json.Unmarshal(jsonData, &conf.Client); err != nil {
		logf("json.Unmarshal(%s, conf) error(%v)", jsonData, err)
		return err
	}

	// 实现监听
	listeners, err := getListenersByFD(conf.Client.Sockets)
	if err != nil {
		logf("getListenersByFD(%v) error(%v)", conf.Client.Sockets, err)
		return err
	}

	ws.initCallback(conf.Client.Config, listeners)
	defer ws.quitCallback()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(signals)

loop:
	for {
		sig := <-signals
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			break loop
		}
	}
	return nil
}
