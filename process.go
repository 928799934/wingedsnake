package wingedsnake

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// fixProcess 发现进程退出就重新启动进程
func fixProcess(ws *wingedSnake, affinity []string, env []string, files []*os.File) {
	ws.wg.Add(1)
	go func() {
		defer ws.wg.Done()

		// 检测间隔
		waitTime := 200 * time.Millisecond
		t := time.NewTimer(waitTime)
		defer t.Stop()
	loop:
		for {
			select {
			case <-ws.closeEvent:
				break loop
			case <-t.C:
			}
			for i, v := range ws.list {
				if v != nil {
					continue
				}
				process, err := startProcess(append(env, fmt.Sprintf("%s=%s", affinityKey, affinity[i])), files)
				if err != nil {
					continue
				}
				ws.list[i] = process
				waitProcess(ws, process)
			}
			t.Reset(waitTime)
		}
	}()
}

func startProcess(env []string, files []*os.File) (*os.Process, error) {
	// Fork exec child process
	name := filepath.Base(os.Args[0]) + " worker process"
	process, err := os.StartProcess(os.Args[0], []string{name}, &os.ProcAttr{Env: env, Files: files})
	if err != nil {
		logf("Fail to fork exec %v", err)
		return nil, err
	}
	return process, nil
}

func waitProcess(ws *wingedSnake, p *os.Process) {
	ws.wg.Add(1)
	go func() {
		defer ws.wg.Done()

		defer func() {
			for i, v := range ws.list {
				if v == nil || v.Pid != p.Pid {
					continue
				}
				ws.list[i] = nil
				v.Release()
				break
			}
		}()

		// 等待进程启动
		time.Sleep(100 * time.Millisecond)

		for {
			state, err := p.Wait()
			if err != nil {
				logf("p.Wait() error(%v)", err)
				return
			}
			if state.Exited() {
				break
			}
		}
	}()
}

func createPidFile(ws *wingedSnake, pidPath string) error {
	szPath := []byte(pidPath)
	if szPath[len(szPath)-1] != '/' {
		pidPath = pidPath + "/"
	}
	_, file := filepath.Split(os.Args[0])
	ws.pidFilePath = pidPath + file + ".pid"
	ws.iPID = os.Getpid()
	strPID := strconv.Itoa(ws.iPID)
	if err := ioutil.WriteFile(ws.pidFilePath, []byte(strPID), 0644); err != nil {
		logf("ioutil.WriteFile(%v,%v,0644) error(%v)", ws.pidFilePath, strPID, err)
		return err
	}
	return nil
}

func removePidFile(ws *wingedSnake) error {
	buf, err := ioutil.ReadFile(ws.pidFilePath)
	if err != nil {
		logf("ioutil.ReadFile(%v) error(%v)", ws.pidFilePath, err)
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
	if err := os.Remove(ws.pidFilePath); err != nil {
		logf("os.Remove(%v) error(%v)", ws.pidFilePath, err)
		return err
	}
	return nil
}
