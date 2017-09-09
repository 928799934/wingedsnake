package main

import (
	"context"
	log "github.com/928799934/log4go.v1"
	ws "github.com/928799934/wingedsnake"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {
	if err := ws.Main(InitFunc, QuitFunc); err != nil {
		panic(err)
	}
}

var (
	ss []*http.Server
)

// InitFunc 测试
func InitFunc(config string, socket, ping, pprof []net.Listener) {
	log.LoadConfiguration(config)

	http.HandleFunc("/ddd", func(wr http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		pid := os.Getpid()
		log.Info("pid[%v] params:%v", pid, r.Form)
	})

	for _, l := range socket {
		s := &http.Server{
			ReadTimeout:    30 * time.Second,
			WriteTimeout:   30 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		ss = append(ss, s)
		go func(s *http.Server, l net.Listener) {
			if err := s.Serve(l); err != nil {
				if err == http.ErrServerClosed {
					return
				}
				log.Error("s.Serve(l) error(%v)", err)
			}
		}(s, l)
	}
	log.Info("init")
}

// QuitFunc 测试
func QuitFunc() {
	log.Info("quit")
	for _, s := range ss {
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		s.Shutdown(ctx)
	}
	log.Close()
}
