// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package wingedsnake

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strconv"
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

	// 获取用户名ID 与用户组ID
	uid, gid, err := getConfigUser(conf)
	if err != nil {
		logf("getConfigUser(%v) error(%v)", conf, err)
		return err
	}

	// 获取CPU亲和数据 与进程数
	affinities := make([]int, len(conf.Base.Affinity))
	for i, v := range conf.Base.Affinity {
		cpuMask, err := strconv.ParseInt(v, 2, 0)
		if err != nil {
			logf("strconv.ParseInt(%v, 2, 0) error(%v)", v, err)
			return err
		}
		affinityMask := 0
		for cpuMask > 0 {
			cpuMask >>= 1
			affinityMask++
		}
		affinities[i] = affinityMask
	}

	jsonData, err := json.Marshal(conf.Client)
	if err != nil {
		logf("json.Marshal(%v) error(%v)", conf, err)
		return err
	}
	// 将配置文件写入子进程环境变量
	env := append(os.Environ(), fmt.Sprintf("%s=%s", globalKey, jsonData))

	// 实现监听
	listeners, err := getListeners(conf.Client.Sockets)
	if err != nil {
		logf("getListeners(%v) error(%v)", conf.Client.Sockets, err)
		return err
	}

	// 装填fd
	files := []*os.File{os.Stdin, os.Stdout, os.Stderr}
	for _, v := range listeners {
		if f := getFileByListener(v); f != nil {
			files = append(files, f)
			continue
		}
		logf("getFileByListener(%v) fail ", v.Addr().String())
		return fmt.Errorf("getFileByListener(%v) fail ", v.Addr().String())
	}

	ws.running = make([]*os.Process, len(affinities))
	// 启动子线程
	thread := newThread(env, files, affinities, ws.running, uid, gid)

	createPidFile(ws, conf.Base.PidPath)
	defer removePidFile(ws, conf.Base.PidPath)

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
			for _, v := range ws.running {
				if v != nil {
					v.Signal(syscall.SIGTERM)
				}
			}
		}
	}
	thread.Close()
	// 等待 thread 退出
	time.Sleep(200 * time.Millisecond)
	// 发送进程退出信号
	for _, v := range ws.running {
		if v != nil {
			v.Signal(syscall.SIGTERM)
		}
	}
	// 等待所有协程 退出
	thread.Wait()

	// 关闭监听
	for _, v := range listeners {
		v.Close()
	}
	// unix socket 需要关闭时间  立即退出会关闭失败
	time.Sleep(100 * time.Millisecond)
	return nil
}
