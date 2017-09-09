package wingedsnake

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// master 主进程
func master(ws *wingedSnake) error {
	confPath := os.Args[1]

	// 读取配置文件
	conf, err := newConfig(confPath)
	if err != nil {
		logf("newConfig(%v) error(%v)", confPath, err)
		return err
	}

	// 实现监听
	pprof, socket, ping, err := getListeners(conf.Client.Pprof, conf.Client.Socket, conf.Client.Ping)
	if err != nil {
		logf("getListeners(%v,%v,%v) error(%v)", conf.Client.Pprof, conf.Client.Socket, conf.Client.Ping, err)
		return err
	}

	// 关闭监听
	defer func() {
		for _, v := range pprof {
			v.Close()
		}
		for _, v := range socket {
			v.Close()
		}
		for _, v := range ping {
			v.Close()
		}
	}()

	jsonData, err := json.Marshal(conf)
	if err != nil {
		logf("json.Marshal(%v) error(%v)", conf, err)
		return err
	}

	// 将配置文件写入子进程环境变量
	env := append(os.Environ(), fmt.Sprintf("%s=%s", globalKey, jsonData))
	// 装填fd
	files := []*os.File{os.Stdin, os.Stdout, os.Stderr}
	for _, v := range pprof {
		f := getFileByListener(v)
		if f == nil {
			logf("pprof getFileByListener(%v) error(%v)", v.Addr().String(), err)
			return err
		}
		files = append(files, f)
	}
	for _, v := range socket {
		f := getFileByListener(v)
		if f == nil {
			logf("socket getFileByListener(%v) error(%v)", v.Addr().String(), err)
			return err
		}
		files = append(files, f)
	}
	for _, v := range ping {
		f := getFileByListener(v)
		if f == nil {
			logf("ping getFileByListener(%v) error(%v)", v.Addr().String(), err)
			return err
		}
		files = append(files, f)
	}

	ws.list = make([]*os.Process, len(conf.Base.Affinity))
	for i, v := range conf.Base.Affinity {
		process, err := startProcess(append(env, fmt.Sprintf("%s=%s", affinityKey, v)), files)
		if err != nil {
			continue
		}
		// 从0开始
		ws.list[i] = process
		waitProcess(ws, process)
	}

	fixProcess(ws, conf.Base.Affinity, env, files)

	createPidFile(ws, conf.Base.PidPath)
	defer removePidFile(ws)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(signals)

loop:
	for {
		sig := <-signals
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			break loop
		case syscall.SIGHUP:
			// 发送进程退出信号
			for _, v := range ws.list {
				if v != nil {
					v.Signal(syscall.SIGTERM)
					v.Release()
				}
			}
		}
	}

	close(ws.closeEvent)
	// 等待 fixProcess 退出
	time.Sleep(100 * time.Millisecond)
	// 发送进程退出信号
	for _, v := range ws.list {
		if v != nil {
			v.Signal(syscall.SIGTERM)
			v.Release()
		}
	}
	// 等待所有协程 退出
	ws.wg.Wait()
	return nil
}
