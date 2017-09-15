package main

import (
	"context"
	log "github.com/928799934/log4go.v1"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

var (
	ss []*http.Server
)

func handleCatFile(wr http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	file := r.FormValue("file")
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		log.Error("ioutil.ReadFile(%v) error(%v)", file, err)
		wr.Write([]byte(err.Error()))
		return
	}
	wr.Write(buf)
}

func startHTTPListen(listeners []net.Listener) {
	http.HandleFunc("/ddd", handleCatFile)

	for _, l := range listeners {
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
}

func stopHTTPListen() {
	for _, s := range ss {
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		s.Shutdown(ctx)
	}
}
