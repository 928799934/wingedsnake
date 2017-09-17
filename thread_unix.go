// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package wingedsnake

import (
	"context"
	"os"
	"runtime"
	"sync"
	"time"
)

type thread struct {
	uid        int
	gid        int
	env        []string
	files      []*os.File
	affinities []int
	list       []*os.Process
	ctx        context.Context
	Close      context.CancelFunc
	sync.WaitGroup
}

func newThread(env []string, files []*os.File, affinities []int, list []*os.Process, uid, gid int) *thread {
	th := &thread{
		uid:        uid,
		gid:        gid,
		env:        env,
		files:      files,
		affinities: affinities,
		list:       list,
	}
	th.ctx, th.Close = context.WithCancel(context.Background())
	th.Add(1)
	go th.running()
	return th
}

func (th *thread) wait(i int) {
	defer th.Done()
	p := th.list[i]
	defer func() {
		p.Release()
		th.list[i] = nil
	}()
	// 等待进程启动
	time.Sleep(100 * time.Millisecond)
	// 调整cpu 亲和
	if len(th.affinities) > i+1 {
		exchangeAffinity(th.affinities[i], p.Pid)
	}
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
}

func (th *thread) running() {
	defer th.Done()
	// 指定后面的C 调用 再同一个线程上
	runtime.LockOSThread()

	// 设置当前线程身份  方便子进程继承该身份
	if err := exchangeOwner(th.uid, th.gid); err != nil {
		logf("exchangeOwner(%v,%v) error(%v)", th.uid, th.gid, err)
		return
	}
	waitTime := 1 * time.Second
	t := time.NewTimer(waitTime)
	defer t.Stop()
loop:
	for {
		select {
		case <-th.ctx.Done():
			break loop
		case <-t.C:
		}
		for i, v := range th.list {
			if v != nil {
				continue
			}
			process, err := startProcess(th.env, th.files)
			if err != nil {
				logf("startProcess(env, files) error(%v)", err)
				continue
			}
			th.list[i] = process
			th.Add(1)
			go th.wait(i)
		}
		t.Reset(waitTime)
	}
}
