package wingedsnake

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

// worker 工作进程函数
func worker(ws *wingedSnake) error {
	// 读取环境变量中的配置
	jsonData := []byte(os.Getenv(globalKey))
	if len(jsonData) == 0 {
		return fmt.Errorf("lost config path")
	}

	// 读取CPU亲和
	affinityMask, err := strconv.ParseInt(os.Getenv(affinityKey), 2, 0)
	if err != nil {
		return fmt.Errorf("config affinity error")
	}

	conf := &config{}
	if err := json.Unmarshal(jsonData, conf); err != nil {
		logf("json.Unmarshal(%s, conf) error(%v)", jsonData, err)
		return err
	}

	// 修改CPU亲和
	if err := exchangeAffinity(int(affinityMask)); err != nil {
		logf("exchangeAffinity(%v) error(%v)", affinityMask, err)
		return err
	}

	// 变换进程uid gid
	configUID, configGID, err := getConfigUser(conf)
	if err != nil {
		logf("getConfigUser(%v) error(%v)", conf, err)
		return err
	}

	if err := exchangeOwner(configUID, configGID); err != nil {
		logf("exchangeOwner(%v,%v) error(%v)", configUID, configGID, err)
		return err
	}

	// 实现监听
	pprof, socket, ping, err := getListenersByFD(conf.Client.Pprof, conf.Client.Socket, conf.Client.Ping)
	if err != nil {
		logf("getListenersByFD(%v,%v,%v) error(%v)", conf.Client.Pprof, conf.Client.Socket, conf.Client.Ping, err)
		return err
	}

	ws.initCallback(conf.Client.Config, socket, ping, pprof)
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
